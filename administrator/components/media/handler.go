package media

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/andycai/weapi/conf"
	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/model"
)

var ErrUnauthorized = errors.New("unauthorized")
var ErrDraftIsInvalid = errors.New("draft is invalid")
var ErrPageIsNotPublish = errors.New("page is not publish")
var ErrPostIsNotPublish = errors.New("post is not publish")
var ErrInvalidPathAndName = errors.New("invalid path and name")
var ErrUploadsDirNotConfigured = errors.New("uploads dir not configured")

func RemoveFile(path, name string) error {
	if name == "" {
		return ErrInvalidPathAndName
	}

	media, err := GetMedia(path, name)
	if err != nil {
		return err
	}

	if !media.External {
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
