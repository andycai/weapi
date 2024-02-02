package model

import (
	"path/filepath"

	"github.com/andycai/weapi/enum"
)

type Media struct {
	BaseContent
	Size       int64  `json:"size"`
	Directory  bool   `json:"directory" gorm:"index"`
	Path       string `json:"path" gorm:"size:200;uniqueIndex:,composite:_path_name"`
	Name       string `json:"name" gorm:"size:200;uniqueIndex:,composite:_path_name"`
	Ext        string `json:"ext" gorm:"size:100"`
	Dimensions string `json:"dimensions" gorm:"size:200"` // x*y
	StorePath  string `json:"-" gorm:"size:300"`
	External   bool   `json:"external"`
	PublicUrl  string `json:"public_url,omitempty" gorm:"-"`
}
type MediaFolder struct {
	Name         string `json:"name"`
	Path         string `json:"path"`
	FilesCount   int64  `json:"filesCount"`
	FoldersCount int64  `json:"foldersCount"`
}

func (m *Media) BuildPublicUrls(mediaHost string, mediaPrefix string) {
	if m.Directory {
		m.PublicUrl = ""
		return
	}

	publicUrl := filepath.Join(mediaPrefix, m.Path, m.Name)
	if mediaHost != "" {
		if mediaHost[len(mediaHost)-1] == '/' {
			mediaHost = mediaHost[:len(mediaHost)-1]
		}
		publicUrl = mediaHost + publicUrl
	}
	m.PublicUrl = publicUrl

	if m.ContentType == enum.ContentTypeImage && m.Thumbnail == "" {
		m.Thumbnail = m.PublicUrl
	}
}
