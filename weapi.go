package weapi

import (
	"embed"
	"path/filepath"

	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/log"
	"github.com/andycai/weapi/model"
)

//go:embed static
var EmbedAssets embed.FS

var ContentTypes = []model.AdminSelectOption{
	{Value: enum.ContentTypeJson, Label: "JSON"},
	{Value: enum.ContentTypeJson, Label: "HTML"},
	{Value: enum.ContentTypeText, Label: "PlainText"},
	{Value: enum.ContentTypeMarkdown, Label: "Markdown"},
	{Value: enum.ContentTypeImage, Label: "Image"},
	{Value: enum.ContentTypeVideo, Label: "Video"},
	{Value: enum.ContentTypeAudio, Label: "Audio"},
	{Value: enum.ContentTypeFile, Label: "File"},
}

var EnabledPageContentTypes = []model.AdminSelectOption{
	{Value: enum.ContentTypeJson, Label: "JSON"},
	{Value: enum.ContentTypeHtml, Label: "HTML"},
	{Value: enum.ContentTypeMarkdown, Label: "Markdown"},
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
