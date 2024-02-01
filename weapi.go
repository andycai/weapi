package weapi

import (
	"embed"
	"path/filepath"

	"github.com/andycai/weapi/enum"
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

func ReadIcon(name string) *model.AdminIcon {
	path := filepath.Join("static/admin/", name)
	data, err := EmbedAssets.ReadFile(path)
	if err != nil {
		// carrot.Warning("Read icon failed:", name, err)
		return nil
	}
	return &model.AdminIcon{SVG: string(data)}
}
