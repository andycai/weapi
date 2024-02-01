package media

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/andycai/weapi/conf"
	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/log"
	"github.com/andycai/weapi/model"
	"github.com/andycai/weapi/utils/random"
	"gorm.io/gorm/clause"
)

func CreateFolder(parent, name string, user *model.User) (string, error) {
	if parent == "" {
		parent = "/"
	}
	obj := model.Media{
		Path:      parent,
		Name:      name,
		Directory: true,
	}

	if user != nil {
		obj.Creator = *user
		obj.CreatorID = user.ID
	}

	fullPath := filepath.Join(parent, name)
	return fullPath, db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&obj).Error
}

func ListFolders(path string) ([]model.MediaFolder, error) {
	var folders []model.MediaFolder = make([]model.MediaFolder, 0)
	tx := db.Model(&model.Media{}).Select("path", "name").Where("path", path).Where("directory", true)
	r := tx.Find(&folders)
	if r.Error != nil {
		return nil, r.Error
	}
	for i := range folders {
		folder := &folders[i]
		folder.Path = filepath.Join(folder.Path, folder.Name)
		tx := db.Model(&model.Media{}).Where("path", folder.Path)
		tx.Select("COUNT(*)").Where("directory", true).Find(&folder.FoldersCount)
		tx = db.Model(&model.Media{}).Where("path", folder.Path)
		tx.Select("COUNT(*)").Where("directory", false).Find(&folder.FilesCount)
	}
	return folders, r.Error
}

func RemoveDirectory(path string) (string, error) {
	var files []model.Media
	r := db.Model(&model.Media{}).Where("path", path).Find(&files)
	if r.Error != nil {
		log.Infof("Remove directory failed: %v, %s", r.Error, path)
		return "", r.Error
	}

	uploadDir := conf.GetValue(db, enum.KEY_CMS_UPLOAD_DIR)
	for _, media := range files {
		if media.Directory {
			RemoveDirectory(filepath.Join(path, media.Name))
			continue
		}
		if !media.External {
			fullPath := filepath.Join(uploadDir, media.StorePath)
			if err := os.Remove(fullPath); err != nil {
				// carrot.Warning("Remove file failed: ", err, fullPath)
			}
		}
	}

	r = db.Where("path", path).Delete(&model.Media{})
	if r.Error != nil {
		return "", r.Error
	}

	parent, name := filepath.Split(path)
	if parent != "/" {
		parent = strings.TrimSuffix(parent, "/")
	}
	return parent, db.Where("path", parent).Where("name", name).Delete(&model.Media{}).Error
}

func RemoveFile(path, name string) error {
	if name == "" {
		return enum.ErrInvalidPathAndName
	}

	media, err := GetMedia(path, name)
	if err != nil {
		return err
	}

	if media.External {
		return nil
	}

	uploadDir := conf.GetValue(db, enum.KEY_CMS_UPLOAD_DIR)
	fullPath := filepath.Join(uploadDir, media.StorePath)
	if err := os.Remove(fullPath); err != nil {
		return err
	}
	return nil
}

func GetMedia(path, name string) (*model.Media, error) {
	var obj model.Media
	if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	tx := db.Model(&model.Media{}).Where("path", path).Where("name", name)
	r := tx.First(&obj)
	if r.Error != nil {
		return nil, r.Error
	}
	return &obj, nil
}

func MakeMediaPublish(siteID, path, name string, obj any, publish bool) error {
	tx := db.Model(obj).Where("site_id", siteID).Where("path", path).Where("name", name)
	vals := map[string]any{"published": publish}
	vals["published"] = publish
	return tx.Updates(vals).Error
}

func PrepareStoreLocalDir() (string, error) {
	uploadDir := conf.GetValue(db, enum.KEY_CMS_UPLOAD_DIR)
	if uploadDir == "" {
		return "", enum.ErrUploadsDirNotConfigured
	}

	if _, err := os.Stat(uploadDir); err != nil {
		if os.IsNotExist(err) {
			// carrot.Warning("upload dir not exist, create it: ", uploadDir)
			if err = os.MkdirAll(uploadDir, 0755); err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	return uploadDir, nil
}

func StoreLocal(uploadDir, storePath string, data []byte) (string, error) {
	storePath = filepath.Join(uploadDir, storePath)
	err := os.WriteFile(storePath, data, 0644)
	if err != nil {
		return "", err
	}
	return storePath, nil
}

func StoreExternal(externalUploader, path, name string, data []byte) (string, error) {
	buf := new(bytes.Buffer)
	form := multipart.NewWriter(buf)
	form.WriteField("path", path)
	form.WriteField("name", name)

	fileField, _ := form.CreateFormFile("file", name)
	fileField.Write(data)
	form.Close()

	resp, err := http.Post(externalUploader, form.FormDataContentType(), buf)
	if err != nil {
		// carrot.Warning("upload to external server failed: ", err, externalUploader)
		return "", err
	}

	defer resp.Body.Close()
	respBody := bytes.NewBuffer(nil)
	io.Copy(respBody, resp.Body)
	body := respBody.Bytes()
	if resp.StatusCode != http.StatusOK {
		// carrot.Warning("upload to external server failed: ", resp.StatusCode, externalUploader, string(body))
		return "", fmt.Errorf("upload to external server failed, code:%d %s", resp.StatusCode, string(body))
	}
	var remoteResult model.UploadResult
	json.Unmarshal(body, &remoteResult)
	return remoteResult.StorePath, nil
}

func UploadFile(path, name string, reader io.Reader) (*model.UploadResult, error) {
	var r model.UploadResult
	r.Path = path
	r.Name = name
	r.Ext = strings.ToLower(filepath.Ext(name))

	canGetDimension := false

	switch r.Ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		canGetDimension = true
		fallthrough
	case ".webp", ".svg", ".ico", ".bmp":
		r.ContentType = enum.ContentTypeImage
	case ".mp3", ".wav", ".ogg", ".aac", ".flac":
		r.ContentType = enum.ContentTypeAudio
	case ".mp4", ".webm", ".avi", ".mov", ".wmv", ".mkv":
		r.ContentType = enum.ContentTypeVideo
	default:
		r.ContentType = enum.ContentTypeFile
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	r.Size = int64(len(data))

	externalUploader := conf.GetValue(db, enum.KEY_CMS_EXTERNAL_UPLOADER)
	if externalUploader != "" {
		storePath, err := StoreExternal(externalUploader, path, name, data)
		if err != nil {
			return nil, err
		}
		r.StorePath = storePath
		r.External = true
	} else {
		storePath := fmt.Sprintf("%s%s", random.RandText(10), r.Ext)
		r.StorePath = storePath
		r.External = false
		uploadDir, err := PrepareStoreLocalDir()
		if err != nil {
			return nil, err
		}
		_, err = StoreLocal(uploadDir, storePath, data)
		if err != nil {
			return nil, err
		}
	}

	if canGetDimension {
		config, _, err := image.DecodeConfig(bytes.NewReader(data))
		if err == nil {
			r.Dimensions = fmt.Sprintf("%dX%d", config.Width, config.Height)
		} else {
			// carrot.Warning("decode image config error: ", err)
			r.Dimensions = "X"
		}
	}
	return &r, nil
}