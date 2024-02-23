package content

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/andycai/weapi/admin/user"
	"github.com/andycai/weapi/constant"
	"github.com/andycai/weapi/log"
	"github.com/andycai/weapi/model"
	"github.com/andycai/weapi/utils/random"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//#region category

func queryCategoryWithCount(siteId, contentObject string) ([]model.Category, error) {
	var tx *gorm.DB
	switch contentObject {
	case "post":
		tx = db.Model(&model.Post{}).Where("site_id", siteId)
	case "page":
		tx = db.Model(&model.Page{}).Where("site_id", siteId)
	default:
		return nil, fmt.Errorf("invalid content object: %s", contentObject)
	}

	var vals []model.Category
	r := db.Model(&model.Category{}).Where("site_id", siteId).Find(&vals)
	if r.Error != nil {
		return nil, r.Error
	}
	for i := range vals {
		val := &vals[i]
		tx := tx.Where("category_id", val.UUID)
		var count int64
		tx.Count(&count)
		val.Count = int(count)
	}
	return vals, r.Error
}

func newRenderCategory(categoryID, categoryPath string) *model.RenderCategory {
	var category model.Category
	r := db.Model(&model.Category{}).Where("uuid", categoryID).First(&category)
	if r.Error != nil {
		return nil
	}

	selected := category.FindItem(categoryPath)

	obj := &model.RenderCategory{
		UUID: category.UUID,
		Name: category.Name,
	}
	if selected != nil {
		obj.Path = selected.Path
		obj.PathName = selected.Name
	}
	return obj
}

// Query tags by category
func GetTagsByCategory(contentType string, form *model.TagsForm) ([]string, error) {
	var tx *gorm.DB
	switch contentType {
	case "post":
		tx = db.Model(&model.Post{})
	case "page":
		tx = db.Model(&model.Page{})
	}

	if form.SiteId != "" {
		tx = tx.Where("site_id", form.SiteId)
	}

	if form.CategoryId != "" {
		tx = tx.Where("category_id", form.CategoryId)
	}

	if form.CategoryPath != "" {
		tx = tx.Where("category_path", form.CategoryPath)
	}

	var rawTags []string
	r := tx.Pluck("tags", &rawTags)
	if r.Error != nil {
		return nil, r.Error
	}

	var uniqueTags map[string]string = make(map[string]string)
	for _, tag := range rawTags {
		if tag == "" {
			continue
		}
		vals := strings.Split(tag, ",")
		for _, val := range vals {
			val = strings.TrimSpace(val)
			if val == "" {
				continue
			}
			uniqueTags[strings.ToLower(val)] = val
		}
	}

	var tags []string = make([]string, 0, len(uniqueTags))
	for k, v := range uniqueTags {
		if k == "" {
			continue
		}
		tags = append(tags, v)
	}
	return tags, r.Error
}

func queryContentByTags(contentType string, form *model.QueryByTagsForm) ([]string, error) {
	return nil, nil
}

//#endregion

//#region page

func NewRenderContentFromPage(page *model.Page) *model.RenderContent {
	var data any
	if page.ContentType == constant.ContentTypeJson {
		data = make(map[string]any)
		err := json.Unmarshal([]byte(page.Body), &data)
		if err != nil {
			log.Infof("unmarshal json error: %s, %s, %s, %s", page.SiteID, page.ID, page.Title, err)
		}
	} else {
		data = page.Body
	}

	return &model.RenderContent{
		BaseContent: page.BaseContent,
		ID:          page.ID,
		SiteID:      page.SiteID,
		PageData:    data,
		IsDraft:     page.IsDraft,
	}
}

func makeDuplicate(obj any) error {
	if page, ok := obj.(*model.Page); ok {
		page.ID = page.ID + "-copy-" + random.RandText(3)
		page.Title = page.Title + "-copy"
		page.IsDraft = true
		page.PreviewURL = ""
		page.Published = false
		page.CreatedAt = time.Now()
		page.UpdatedAt = time.Now()
		return db.Create(page).Error
	} else if post, ok := obj.(*model.Post); ok {
		post.ID = post.ID + "-copy-" + random.RandText(3)
		post.Title = post.Title + "-copy"
		post.IsDraft = true
		post.PreviewURL = ""
		post.CreatedAt = time.Now()
		post.UpdatedAt = time.Now()
		post.Published = false
		return db.Create(post).Error
	}
	return errors.New("invalid object, must be page or post")
}

func makePublish(siteID, ID string, obj any, publish bool) error {
	tx := db.Model(obj).Where("site_id", siteID).Where("id", ID)
	vals := map[string]any{"published": publish}

	vals["published"] = publish
	if publish {
		vals["body"] = gorm.Expr("draft")
		vals["is_draft"] = false
	}
	return tx.Updates(vals).Error
}

func safeDraft(siteID, ID string, obj any, draft string) error {
	tx := db.Model(obj).Where("site_id", siteID).Where("id", ID)
	vals := map[string]any{
		"is_draft": true,
		"draft":    draft,
	}
	return tx.Updates(vals).Error
}

func queryPageTags() ([]string, error) {
	var vals []string
	r := db.Model(&model.Page{}).Select("DISTINCT(tags)").Find(&vals)
	return vals, r.Error
}

//#endregion

//#region post

func NewRenderContentFromPost(post *model.Post, relations bool) *model.RenderContent {
	r := &model.RenderContent{
		BaseContent: post.BaseContent,
		ID:          post.ID,
		SiteID:      post.SiteID,
		PostBody:    post.Body,
		IsDraft:     post.IsDraft,
		Category:    newRenderCategory(post.CategoryID, post.CategoryPath),
	}

	if relations {
		relationCount := user.GetIntValue(constant.KEY_CMS_RELATION_COUNT, 3)
		suggestionCount := user.GetIntValue(constant.KEY_CMS_SUGGESTION_COUNT, 3)

		r.Relations, _ = GetRelations(post.SiteID, post.CategoryID, post.CategoryPath, post.ID, relationCount)
		r.Suggestions, _ = GetSuggestions(post.SiteID, post.CategoryID, post.CategoryPath, post.ID, suggestionCount)
	}
	return r
}

func queryPostTags() ([]string, error) {
	var vals []string
	r := db.Model(&model.Post{}).Select("DISTINCT(tags)").Find(&vals)
	return vals, r.Error
}

func GetSuggestions(siteId, categoryId, categoryPath, postId string, maxCount int) ([]model.RelationContent, error) {
	return GetRelations(siteId, categoryId, categoryPath, postId, maxCount)
}

func GetRelations(siteId, categoryId, categoryPath, postId string, maxCount int) ([]model.RelationContent, error) {
	if maxCount <= 0 {
		return nil, nil
	}
	var r []model.RelationContent
	tx := db.Model(&model.Post{}).Where("site_id", siteId).Where("published", true)
	if categoryId != "" {
		tx = tx.Where("category_id", categoryId)
	}

	var totalCount int64
	tx.Count(&totalCount)
	if totalCount == 0 {
		return nil, nil
	}
	excludeIds := []string{}
	if postId != "" {
		excludeIds = append(excludeIds, postId)
	}
	for i := 0; i < maxCount; i++ {
		// random select
		offset := rand.Intn(int(totalCount))
		var val model.Post
		subTx := tx
		if len(excludeIds) > 0 {
			subTx = subTx.Where("id NOT IN (?)", excludeIds)
		}
		result := subTx.Offset(offset).Limit(1).Take(&val)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				continue
			}
			return nil, result.Error
		}

		excludeIds = append(excludeIds, val.ID)

		r = append(r, model.RelationContent{
			BaseContent: val.BaseContent,
			SiteID:      val.SiteID,
			ID:          val.ID})
	}
	return r, nil
}

//#endregion

//#region media

func createFolder(parent, name string, user *model.User) (string, error) {
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

func listFolders(path string) ([]model.MediaFolder, error) {
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

func removeDirectory(path string) (string, error) {
	var files []model.Media
	r := db.Model(&model.Media{}).Where("path", path).Find(&files)
	if r.Error != nil {
		log.Infof("Remove directory failed: %v, %s", r.Error, path)
		return "", r.Error
	}

	uploadDir := user.GetValue(constant.KEY_CMS_UPLOAD_DIR)
	for _, media := range files {
		if media.Directory {
			removeDirectory(filepath.Join(path, media.Name))
			continue
		}
		if !media.External {
			fullPath := filepath.Join(uploadDir, media.StorePath)
			if err := os.Remove(fullPath); err != nil {
				log.Infof("Remove file failed: %s, %s", err, fullPath)
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

func removeFile(path, name string) error {
	if name == "" {
		return constant.ErrInvalidPathAndName
	}

	media, err := getMedia(path, name)
	if err != nil {
		return err
	}

	if media.External {
		return nil
	}

	uploadDir := user.GetValue(constant.KEY_CMS_UPLOAD_DIR)
	fullPath := filepath.Join(uploadDir, media.StorePath)
	if err := os.Remove(fullPath); err != nil {
		return err
	}
	return nil
}

func getMedia(path, name string) (*model.Media, error) {
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

func makeMediaPublish(siteID, path, name string, obj any, publish bool) error {
	tx := db.Model(obj).Where("site_id", siteID).Where("path", path).Where("name", name)
	vals := map[string]any{"published": publish}
	vals["published"] = publish
	return tx.Updates(vals).Error
}

func prepareStoreLocalDir() (string, error) {
	uploadDir := user.GetValue(constant.KEY_CMS_UPLOAD_DIR)
	if uploadDir == "" {
		return "", constant.ErrUploadsDirNotConfigured
	}

	if _, err := os.Stat(uploadDir); err != nil {
		if os.IsNotExist(err) {
			log.Infof("upload dir not exist, create it: %s", uploadDir)
			if err = os.MkdirAll(uploadDir, 0755); err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	return uploadDir, nil
}

func storeLocal(uploadDir, storePath string, data []byte) (string, error) {
	storePath = filepath.Join(uploadDir, storePath)
	err := os.WriteFile(storePath, data, 0644)
	if err != nil {
		return "", err
	}
	return storePath, nil
}

func storeExternal(externalUploader, path, name string, data []byte) (string, error) {
	buf := new(bytes.Buffer)
	form := multipart.NewWriter(buf)
	form.WriteField("path", path)
	form.WriteField("name", name)

	fileField, _ := form.CreateFormFile("file", name)
	fileField.Write(data)
	form.Close()

	resp, err := http.Post(externalUploader, form.FormDataContentType(), buf)
	if err != nil {
		log.Infof("upload to external server failed: %s, %s", err, externalUploader)
		return "", err
	}

	defer resp.Body.Close()
	respBody := bytes.NewBuffer(nil)
	io.Copy(respBody, resp.Body)
	body := respBody.Bytes()
	if resp.StatusCode != http.StatusOK {
		log.Infof("upload to external server failed: %s, %s, %s", resp.StatusCode, externalUploader, string(body))
		return "", fmt.Errorf("upload to external server failed, code:%d %s", resp.StatusCode, string(body))
	}
	var remoteResult model.UploadResult
	json.Unmarshal(body, &remoteResult)
	return remoteResult.StorePath, nil
}

func uploadFile(path, name string, reader io.Reader) (*model.UploadResult, error) {
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
		r.ContentType = constant.ContentTypeImage
	case ".mp3", ".wav", ".ogg", ".aac", ".flac":
		r.ContentType = constant.ContentTypeAudio
	case ".mp4", ".webm", ".avi", ".mov", ".wmv", ".mkv":
		r.ContentType = constant.ContentTypeVideo
	default:
		r.ContentType = constant.ContentTypeFile
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	r.Size = int64(len(data))

	externalUploader := user.GetValue(constant.KEY_CMS_EXTERNAL_UPLOADER)
	if externalUploader != "" {
		storePath, err := storeExternal(externalUploader, path, name, data)
		if err != nil {
			return nil, err
		}
		r.StorePath = storePath
		r.External = true
	} else {
		storePath := fmt.Sprintf("%s%s", random.RandText(10), r.Ext)
		r.StorePath = storePath
		r.External = false
		uploadDir, err := prepareStoreLocalDir()
		if err != nil {
			return nil, err
		}
		_, err = storeLocal(uploadDir, storePath, data)
		if err != nil {
			return nil, err
		}
	}

	if canGetDimension {
		config, _, err := image.DecodeConfig(bytes.NewReader(data))
		if err == nil {
			r.Dimensions = fmt.Sprintf("%dX%d", config.Width, config.Height)
		} else {
			log.Infof("decode image config error: %s", err)
			r.Dimensions = "X"
		}
	}
	return &r, nil
}

//#endregion
