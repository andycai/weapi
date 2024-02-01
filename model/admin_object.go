package model

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

type AdminBuildContext func(*fiber.Ctx, map[string]any) map[string]any

type AdminQueryResult struct {
	TotalCount int              `json:"total,omitempty"`
	Pos        int              `json:"pos,omitempty"`
	Limit      int              `json:"limit,omitempty"`
	Keyword    string           `json:"keyword,omitempty"`
	Items      []map[string]any `json:"items"`
	Objects    []any            `json:"-"`
}

// Access control
type AdminAccessCheck func(c *fiber.Ctx, obj *AdminObject) error
type AdminActionHandler func(c *fiber.Ctx, obj any) (any, error)

type AdminSelectOption struct {
	Label string `json:"label"`
	Value any    `json:"value"`
}

type AdminAttribute struct {
	Default      any                 `json:"default,omitempty"`
	Choices      []AdminSelectOption `json:"choices,omitempty"`
	SingleChoice bool                `json:"singleChoice,omitempty"`
	Widget       string              `json:"widget,omitempty"`
	FilterWidget string              `json:"filterWidget,omitempty"`
	Help         string              `json:"help,omitempty"`
}
type AdminForeign struct {
	Path       string `json:"path"`
	Field      string `json:"field"`
	FieldName  string `json:"-"`
	ForeignKey string `json:"-"`
}
type AdminValue struct {
	Value any    `json:"value"`
	Label string `json:"label,omitempty"`
}
type AdminIcon struct {
	Url string `json:"url,omitempty"`
	SVG string `json:"svg,omitempty"`
}

type AdminField struct {
	Placeholder string          `json:"placeholder,omitempty"` // Placeholder of the filed
	Label       string          `json:"label"`                 // Label of the filed
	NotColumn   bool            `json:"notColumn,omitempty"`   // Not a column
	Required    bool            `json:"required,omitempty"`
	Name        string          `json:"name"`
	Type        string          `json:"type"`
	Tag         string          `json:"tag,omitempty"`
	Attribute   *AdminAttribute `json:"attribute,omitempty"`
	CanNull     bool            `json:"canNull,omitempty"`
	IsArray     bool            `json:"isArray,omitempty"`
	Primary     bool            `json:"primary,omitempty"`
	Foreign     *AdminForeign   `json:"foreign,omitempty"`
	IsAutoID    bool            `json:"isAutoId,omitempty"`
	IsPtr       bool            `json:"isPtr,omitempty"`
	ElemType    reflect.Type    `json:"-"`
	FieldName   string          `json:"-"`
}
type AdminScript struct {
	Src    string `json:"src"`
	Onload bool   `json:"onload,omitempty"`
}
type AdminAction struct {
	Path          string             `json:"path"`
	Name          string             `json:"name"`
	Label         string             `json:"label,omitempty"`
	Icon          string             `json:"icon,omitempty"`
	Class         string             `json:"class,omitempty"`
	WithoutObject bool               `json:"withoutObject"`
	Handler       AdminActionHandler `json:"-"`
}

type AdminObject struct {
	Model       any             `json:"-"`
	Group       string          `json:"group"`               // Group name
	Name        string          `json:"name"`                // Name of the object
	Desc        string          `json:"desc,omitempty"`      // Description
	Path        string          `json:"path"`                // Path prefix
	Shows       []string        `json:"shows"`               // Show fields
	Orders      []Order         `json:"orders"`              // Default orders of the object
	Editables   []string        `json:"editables"`           // Editable fields
	Filterables []string        `json:"filterables"`         // Filterable fields
	Orderables  []string        `json:"orderables"`          // Orderable fields, can override Orders
	Searchables []string        `json:"searchables"`         // Searchable fields
	Requireds   []string        `json:"requireds,omitempty"` // Required fields
	PrimaryKeys []string        `json:"primaryKeys"`         // Primary keys name
	UniqueKeys  []string        `json:"uniqueKeys"`          // Primary keys name
	PluralName  string          `json:"pluralName"`
	Fields      []AdminField    `json:"fields"`
	EditPage    string          `json:"editpage,omitempty"`
	ListPage    string          `json:"listpage,omitempty"`
	Scripts     []AdminScript   `json:"scripts,omitempty"`
	Styles      []string        `json:"styles,omitempty"`
	Permissions map[string]bool `json:"permissions,omitempty"`
	Actions     []AdminAction   `json:"actions,omitempty"`
	Icon        *AdminIcon      `json:"icon,omitempty"`
	Invisible   bool            `json:"invisible,omitempty"`

	Attributes       map[string]AdminAttribute `json:"-"` // Field's extra attributes
	AccessCheck      AdminAccessCheck          `json:"-"` // Access control function
	GetDB            GetDB                     `json:"-"`
	BeforeCreate     BeforeCreateFunc          `json:"-"`
	BeforeRender     BeforeRenderFunc          `json:"-"`
	BeforeUpdate     BeforeUpdateFunc          `json:"-"`
	BeforeDelete     BeforeDeleteFunc          `json:"-"`
	TableName        string                    `json:"-"`
	ModelElem        reflect.Type              `json:"-"`
	Ignores          map[string]bool           `json:"-"`
	PrimaryKeyMaping map[string]string         `json:"-"`
	MarkDeletedField string                    `json:"-"`
	Weight           int                       `json:"-"`
}

type ContentIcon AdminIcon

func (s ContentIcon) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *ContentIcon) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &s)
}

type UploadResult struct {
	PublicUrl   string `json:"publicUrl"`
	Thumbnail   string `json:"thumbnail"`
	Path        string `json:"path"`
	Name        string `json:"name"`
	External    bool   `json:"external"`
	StorePath   string `json:"storePath"`
	Dimensions  string `json:"dimensions"`
	Ext         string `json:"ext"`
	Size        int64  `json:"size"`
	ContentType string `json:"contentType"`
}
