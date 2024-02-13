package weapi

import (
	"embed"
	"path/filepath"

	"github.com/andycai/weapi/constant"
	"github.com/andycai/weapi/log"
	"github.com/andycai/weapi/model"
)

//go:embed static
var EmbedAssets embed.FS

var ContentTypes = []model.AdminSelectOption{
	{Value: constant.ContentTypeJson, Label: "JSON"},
	{Value: constant.ContentTypeJson, Label: "HTML"},
	{Value: constant.ContentTypeText, Label: "PlainText"},
	{Value: constant.ContentTypeMarkdown, Label: "Markdown"},
	{Value: constant.ContentTypeImage, Label: "Image"},
	{Value: constant.ContentTypeVideo, Label: "Video"},
	{Value: constant.ContentTypeAudio, Label: "Audio"},
	{Value: constant.ContentTypeFile, Label: "File"},
}

var EnabledPageContentTypes = []model.AdminSelectOption{
	{Value: constant.ContentTypeJson, Label: "JSON"},
	{Value: constant.ContentTypeHtml, Label: "HTML"},
	{Value: constant.ContentTypeMarkdown, Label: "Markdown"},
}

var models = []any{
	&model.User{},
	&model.Site{},
	&model.Category{},
	&model.Page{},
	&model.Post{},
	&model.Media{},
	&model.Comment{},
	&model.Group{},
	&model.GroupMember{},
	&model.Config{},
	&model.Activity{},
	&model.Club{},
	&model.ClubMember{},
}

func ReadIcon(name string) *model.AdminIcon {
	path := filepath.Join("static/admin/", name)
	data, err := EmbedAssets.ReadFile(path)
	if err != nil {
		log.Infof("Read icon failed: %s, %s", name, err)
		return nil
	}
	return &model.AdminIcon{SVG: string(data)}
}
