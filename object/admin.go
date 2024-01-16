package object

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const KEY_ADMIN_DASHBOARD = "ADMIN_DASHBOARD"

type AdminBuildContext func(*fiber.Ctx, map[string]any) map[string]any

type AdminQueryResult struct {
	TotalCount int              `json:"total,omitempty"`
	Pos        int              `json:"pos,omitempty"`
	Limit      int              `json:"limit,omitempty"`
	Keyword    string           `json:"keyword,omitempty"`
	Items      []map[string]any `json:"items"`
	objects    []any            `json:"-"`
}

// Access control
type AdminAccessCheck func(c *fiber.Ctx, obj *AdminObject) error
type AdminActionHandler func(db *gorm.DB, c *fiber.Ctx, obj any) (any, error)

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
	fieldName  string `json:"-"`
	foreignKey string `json:"-"`
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
	elemType    reflect.Type    `json:"-"`
	fieldName   string          `json:"-"`
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
	tableName        string                    `json:"-"`
	modelElem        reflect.Type              `json:"-"`
	ignores          map[string]bool           `json:"-"`
	primaryKeyMaping map[string]string         `json:"-"`
	markDeletedField string                    `json:"-"`
}

type ContentIcon AdminIcon

func (s ContentIcon) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *ContentIcon) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &s)
}
