package media

import (
	"path/filepath"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/andycai/weapi/administrator/components/config"
	"github.com/andycai/weapi/administrator/components/user"
	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/log"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
)

func handleListFolders(c *fiber.Ctx, obj any) (any, error) {
	path := c.Query("path")
	return ListFolders(path)
}

func handleNewFolder(c *fiber.Ctx, obj any) (any, error) {
	path := c.Query("path")
	name := c.Query("name")
	user := user.CurrentUser(c)
	return CreateFolder(path, name, user)
}

func handleMakeMediaPublish(c *fiber.Ctx, obj any, publish bool) (any, error) {
	siteId := c.Query("site_id")
	path := c.Query("path")
	name := c.Query("name")

	if err := MakeMediaPublish(siteId, path, name, obj, publish); err != nil {
		log.Infof("Make publish failed: %s, %s, %s, %t, %v", siteId, path, name, publish, err)
		return false, err
	}
	return true, nil
}

func handleMedia(c *fiber.Ctx) error {
	fullPath := c.Params("*")
	path, name := filepath.Split(fullPath)
	if len(path) > 1 && path[0] != '/' {
		path = "/" + path
	}
	img, err := GetMedia(path, name)
	if err != nil {
		// carrot.AbortWithJSONError(c, http.StatusNotFound, err)
		return err
	}

	if img.External {
		c.Redirect(img.StorePath)
		return nil
	}

	uploadDir := config.GetValue(enum.KEY_CMS_UPLOAD_DIR)
	filepath := filepath.Join(uploadDir, img.StorePath)
	return c.SendFile(filepath)
	// http.ServeFile(c.Request().BodyWriter(), c.Request(), filepath)
}

func handleRemoveDirectory(c *fiber.Ctx, obj any) (any, error) {
	path := c.Query("path")

	parent, err := RemoveDirectory(path)
	if err != nil {
		// carrot.AbortWithJSONError(c, http.StatusInternalServerError, err)
		return nil, err
	}
	return parent, nil
}

func handleUpload(c *fiber.Ctx, obj any) (any, error) {
	created := c.Query("created")
	path := c.Query("path")
	name := c.Query("name")

	file, err := c.FormFile("file")
	if err != nil {
		return nil, err
	}

	mFile, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer mFile.Close()

	if path == "" {
		path = "/"
	}
	if name == "" {
		name = file.Filename
	}
	r, err := UploadFile(path, name, mFile)
	if err != nil {
		return nil, err
	}

	var media model.Media

	user := user.CurrentUser(c)
	media.Name = r.Name
	media.Path = r.Path
	media.External = r.External
	media.StorePath = r.StorePath
	media.Size = r.Size
	media.ContentType = r.ContentType
	media.Dimensions = r.Dimensions
	media.Directory = false
	media.Ext = r.Ext
	media.ContentType = r.ContentType
	media.Published = true

	if user != nil {
		media.Creator = *user
		media.CreatorID = user.ID
	}

	if created != "" {
		result := db.Create(&media)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	mediaHost := config.GetValue(enum.KEY_CMS_MEDIA_HOST)
	mediaPrefix := config.GetValue(enum.KEY_CMS_MEDIA_PREFIX)
	media.BuildPublicUrls(mediaHost, mediaPrefix)

	r.PublicUrl = media.PublicUrl
	r.Thumbnail = media.Thumbnail

	return r, nil
}
