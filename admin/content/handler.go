package content

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/andycai/weapi/admin/user"
	"github.com/andycai/weapi/constant"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/log"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
)

//#region category

func handleQueryCategoryWithCount(c *fiber.Ctx, obj any) (any, error) {
	siteId := c.Query("site_id")
	current := strings.ToLower(c.Query("current"))
	return queryCategoryWithCount(siteId, current)
}

//#endregion

//#region page

func handleMakePagePublish(c *fiber.Ctx, obj any, publish bool) (any, error) {
	siteId := c.Query("site_id")
	id := c.Query("id")
	if err := makePublish(siteId, id, obj, publish); err != nil {
		log.Infof("make publish failed: %s, %s, %t, %s", siteId, id, publish, err)
		return false, err
	}
	return true, nil
}

func handleMakePageDuplicate(c *fiber.Ctx, obj any) (any, error) {
	if err := makeDuplicate(obj); err != nil {
		log.Infof("make duplicate failed: %v, %s", obj, err)
		return false, err
	}
	return true, nil
}

func handleSaveDraft(c *fiber.Ctx, obj any) (any, error) {
	siteId := c.Query("site_id")
	id := c.Query("id")

	var formData map[string]string
	if err := c.BodyParser(&formData); err != nil {
		return nil, err
	}

	draft, ok := formData["draft"]
	if !ok {
		return nil, constant.ErrDraftIsInvalid
	}

	if err := safeDraft(siteId, id, obj, draft); err != nil {
		log.Infof("safe draft failed: %s, %s, %s", siteId, id, err)
		return false, err
	}
	return true, nil
}

func handleQueryPageTags(c *fiber.Ctx, obj any, tableName string) (any, error) {
	return queryPageTags()
}

//#endregion

//#region post

func handleQueryPostTags(c *fiber.Ctx, obj any, tableName string) (any, error) {
	return queryPostTags()
}

//#endregion

//#region media

func handleListFolders(c *fiber.Ctx, obj any) (any, error) {
	path := c.Query("path")
	return listFolders(path)
}

func handleNewFolder(c *fiber.Ctx, obj any) (any, error) {
	path := c.Query("path")
	name := c.Query("name")
	user := user.Current(c)
	return createFolder(path, name, user)
}

func handleMakeMediaPublish(c *fiber.Ctx, obj any, publish bool) (any, error) {
	siteId := c.Query("site_id")
	path := c.Query("path")
	name := c.Query("name")

	if err := makeMediaPublish(siteId, path, name, obj, publish); err != nil {
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
	img, err := getMedia(path, name)
	if err != nil {
		return core.Error(c, http.StatusNotFound, err)
	}

	if img.External {
		return c.Redirect(img.StorePath)
	}

	uploadDir := user.GetValue(constant.KEY_CMS_UPLOAD_DIR)
	filepath := filepath.Join(uploadDir, img.StorePath)
	return c.SendFile(filepath)
}

func handleRemoveDirectory(c *fiber.Ctx, obj any) (any, error) {
	path := c.Query("path")

	parent, err := removeDirectory(path)
	if err != nil {
		return nil, core.Error(c, http.StatusInternalServerError, err)
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
	r, err := uploadFile(path, name, mFile)
	if err != nil {
		return nil, err
	}

	var media model.Media

	userVo := user.Current(c)
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

	if userVo != nil {
		media.Creator = *userVo
		media.CreatorID = userVo.ID
	}

	if created != "" {
		result := db.Create(&media)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	mediaHost := user.GetValue(constant.KEY_CMS_MEDIA_HOST)
	mediaPrefix := user.GetValue(constant.KEY_CMS_MEDIA_PREFIX)
	media.BuildPublicUrls(mediaHost, mediaPrefix)

	r.PublicUrl = media.PublicUrl
	r.Thumbnail = media.Thumbnail

	return r, nil
}

//#endregion
