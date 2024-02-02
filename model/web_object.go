package model

import (
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

const (
	DefaultQueryLimit = 102400 // 100k
)

const (
	FilterOpIsNot          = "is not"
	FilterOpEqual          = "="
	FilterOpNotEqual       = "<>"
	FilterOpIn             = "in"
	FilterOpNotIn          = "not_in"
	FilterOpGreater        = ">"
	FilterOpGreaterOrEqual = ">="
	FilterOpLess           = "<"
	FilterOpLessOrEqual    = "<="
	FilterOpLike           = "like"
	FilterOpBetween        = "between"
)

const (
	OrderOpDesc = "desc"
	OrderOpAsc  = "asc"
)

const (
	GET    = 1 << 1
	CREATE = 1 << 2
	EDIT   = 1 << 3
	DELETE = 1 << 4
	QUERY  = 1 << 5
)

type PrepareQuery func(c *fiber.Ctx) (*QueryForm, error)

type (
	BeforeCreateFunc      func(ctx *fiber.Ctx, vptr any) error
	BeforeDeleteFunc      func(ctx *fiber.Ctx, vptr any) error
	BeforeUpdateFunc      func(ctx *fiber.Ctx, vptr any, vals map[string]any) error
	BeforeRenderFunc      func(ctx *fiber.Ctx, vptr any) (any, error)
	BeforeQueryRenderFunc func(ctx *fiber.Ctx, r *QueryResult) (any, error)
)

type QueryView struct {
	Path    string `json:"path"`
	Method  string `json:"method"`
	Desc    string `json:"desc"`
	Prepare PrepareQuery
}

type WebObjectPrimaryField struct {
	IsPrimary bool
	Name      string
	Kind      reflect.Kind
	JSONName  string
}

type WebObject struct {
	Model             any
	Group             string
	Name              string
	Desc              string
	AuthRequired      bool
	Editables         []string
	Filterables       []string
	Orderables        []string
	Searchables       []string
	BeforeCreate      BeforeCreateFunc
	BeforeUpdate      BeforeUpdateFunc
	BeforeDelete      BeforeDeleteFunc
	BeforeRender      BeforeRenderFunc
	BeforeQueryRender BeforeQueryRenderFunc

	Views        []QueryView
	AllowMethods int

	PrimaryKeys []WebObjectPrimaryField
	UniqueKeys  []WebObjectPrimaryField
	TableName   string

	// Model type
	ModelElem reflect.Type
	// Map json tag to struct field name. such as:
	// UUID string `json:"id"` => {"id" : "UUID"}
	JsonToFields map[string]string
	// Map json tag to field kind. such as:
	// UUID string `json:"id"` => {"id": string}
	JsonToKinds map[string]reflect.Kind
}

type Filter struct {
	IsTimeType bool   `json:"-"`
	Name       string `json:"name"`
	Op         string `json:"op"`
	Value      any    `json:"value"`
}

type Order struct {
	Name string `json:"name"`
	Op   string `json:"op"`
}

type QueryForm struct {
	Pos          int      `json:"pos"`
	Limit        int      `json:"limit"`
	Keyword      string   `json:"keyword,omitempty"`
	Filters      []Filter `json:"filters,omitempty"`
	Orders       []Order  `json:"orders,omitempty"`
	ForeignMode  bool     `json:"foreign"` // for foreign key
	ViewFields   []string `json:"-"`       // for view
	SearchFields []string `json:"-"`       // for keyword
}

type QueryResult struct {
	TotalCount int    `json:"total,omitempty"`
	Pos        int    `json:"pos,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	Keyword    string `json:"keyword,omitempty"`
	Items      []any  `json:"items"`
}

// GetQuery return the combined filter SQL statement.
// such as "age >= ?", "name IN ?".
func (f *Filter) GetQuery() string {
	var op string
	switch f.Op {
	case FilterOpIsNot:
		op = "IS NOT"
	case FilterOpEqual:
		op = "="
	case FilterOpNotEqual:
		op = "<>"
	case FilterOpIn:
		op = "IN"
	case FilterOpNotIn:
		op = "NOT IN"
	case FilterOpGreater:
		op = ">"
	case FilterOpGreaterOrEqual:
		op = ">="
	case FilterOpLess:
		op = "<"
	case FilterOpLessOrEqual:
		op = "<="
	case FilterOpLike:
		op = "LIKE"
	case FilterOpBetween:
		op = "BETWEEN"
		return fmt.Sprintf("`%s` BETWEEN ? AND ?", f.Name)
	}

	if op == "" {
		return ""
	}

	return fmt.Sprintf("`%s` %s ?", f.Name, op)
}

// GetQuery return the combined order SQL statement.
// such as "id DESC".
func (f *Order) GetQuery() string {
	if f.Op == OrderOpDesc {
		return f.Name + " DESC"
	}
	return f.Name + " ASC"
}
