package model

import (
	"database/sql/driver"
	"encoding/json"
)

type CategoryItem struct {
	Path     string        `json:"path"`
	Name     string        `json:"name"`
	Icon     *ContentIcon  `json:"icon,omitempty"`
	Children CategoryItems `json:"children,omitempty"`
	Count    int           `json:"count" gorm:"-"`
}

type Category struct {
	SiteID string        `json:"site_id" gorm:"uniqueIndex:,composite:_site_uuid"`
	Site   Site          `json:"-"`
	UUID   string        `json:"uuid" gorm:"size:12;uniqueIndex:,composite:_site_uuid"`
	Name   string        `json:"name" gorm:"size:200"`
	Items  CategoryItems `json:"items,omitempty"`
	Count  int           `json:"count" gorm:"-"`
}

type RenderCategory struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Path     string `json:"path,omitempty"`
	PathName string `json:"path_name,omitempty"`
}

func (s CategoryItem) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *CategoryItem) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &s)
}

type CategoryItems []CategoryItem

func (s CategoryItems) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *CategoryItems) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &s)
}

func (category *Category) findItem(path string, items CategoryItems) *CategoryItem {
	for _, item := range items {
		if item.Path == path {
			return &item
		}

		if item.Children != nil {
			if found := category.findItem(path, item.Children); found != nil {
				return found
			}
		}
	}
	return nil
}

func (category *Category) FindItem(path string) *CategoryItem {
	if path == "" {
		return nil
	}
	return category.findItem(path, category.Items)
}
