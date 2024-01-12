package weapi

import (
	"embed"

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
