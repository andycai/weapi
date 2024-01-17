package weapi

import (
	"embed"
	"path/filepath"

	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/object"
)

//go:embed static
var EmbedAssets embed.FS

var ContentTypes = []object.AdminSelectOption{
	{Value: enum.ContentTypeJson, Label: "JSON"},
	{Value: enum.ContentTypeJson, Label: "HTML"},
	{Value: enum.ContentTypeText, Label: "PlainText"},
	{Value: enum.ContentTypeMarkdown, Label: "Markdown"},
	{Value: enum.ContentTypeImage, Label: "Image"},
	{Value: enum.ContentTypeVideo, Label: "Video"},
	{Value: enum.ContentTypeAudio, Label: "Audio"},
	{Value: enum.ContentTypeFile, Label: "File"},
}

var EnabledPageContentTypes = []object.AdminSelectOption{
	{Value: enum.ContentTypeJson, Label: "JSON"},
	{Value: enum.ContentTypeHtml, Label: "HTML"},
	{Value: enum.ContentTypeMarkdown, Label: "Markdown"},
}

func ReadIcon(name string) *object.AdminIcon {
	data, err := EmbedAssets.ReadFile(filepath.Join("admin", name))
	if err != nil {
		// carrot.Warning("Read icon failed:", name, err)
		return nil
	}
	return &object.AdminIcon{SVG: string(data)}
}
