package object

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

type GetDB func(c *fiber.Ctx, isCreate bool) *gorm.DB // designed for group
type PrepareQuery func(db *gorm.DB, c *fiber.Ctx) (*gorm.DB, *QueryForm, error)

type (
	BeforeCreateFunc      func(db *gorm.DB, ctx *fiber.Ctx, vptr any) error
	BeforeDeleteFunc      func(db *gorm.DB, ctx *fiber.Ctx, vptr any) error
	BeforeUpdateFunc      func(db *gorm.DB, ctx *fiber.Ctx, vptr any, vals map[string]any) error
	BeforeRenderFunc      func(db *gorm.DB, ctx *fiber.Ctx, vptr any) (any, error)
	BeforeQueryRenderFunc func(db *gorm.DB, ctx *fiber.Ctx, r *QueryResult) (any, error)
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
	GetDB             GetDB
	BeforeCreate      BeforeCreateFunc
	BeforeUpdate      BeforeUpdateFunc
	BeforeDelete      BeforeDeleteFunc
	BeforeRender      BeforeRenderFunc
	BeforeQueryRender BeforeQueryRenderFunc

	Views        []QueryView
	AllowMethods int

	primaryKeys []WebObjectPrimaryField
	uniqueKeys  []WebObjectPrimaryField
	tableName   string

	// Model type
	modelElem reflect.Type
	// Map json tag to struct field name. such as:
	// UUID string `json:"id"` => {"id" : "UUID"}
	jsonToFields map[string]string
	// Map json tag to field kind. such as:
	// UUID string `json:"id"` => {"id": string}
	jsonToKinds map[string]reflect.Kind
}

type Filter struct {
	isTimeType bool   `json:"-"`
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
	searchFields []string `json:"-"`       // for keyword
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

func (obj *WebObject) RegisterObject(r fiber.Router) error {
	if err := obj.Build(); err != nil {
		return err
	}

	p := obj.Name
	allowMethods := obj.AllowMethods
	if allowMethods == 0 {
		allowMethods = GET | CREATE | EDIT | DELETE | QUERY
	}

	primaryKeyPath := obj.BuildPrimaryPath(p)
	if allowMethods&GET != 0 {
		r.Get(primaryKeyPath, func(c *fiber.Ctx) error {
			// handleGetObject(c, obj)
			return nil
		})
	}
	if allowMethods&CREATE != 0 {
		r.Put(p, func(c *fiber.Ctx) error {
			// handleCreateObject(c, obj)
			return nil
		})
	}
	if allowMethods&EDIT != 0 {
		r.Patch(primaryKeyPath, func(c *fiber.Ctx) error {
			// handleEditObject(c, obj)
			return nil
		})
	}

	if allowMethods&DELETE != 0 {
		r.Delete(primaryKeyPath, func(c *fiber.Ctx) error {
			// handleDeleteObject(c, obj)
			return nil
		})
	}

	if allowMethods&QUERY != 0 {
		r.Post(p, func(c *fiber.Ctx) error {
			// handleQueryObject(c, obj, DefaultPrepareQuery)
			return nil
		})
	}

	for i := 0; i < len(obj.Views); i++ {
		v := &obj.Views[i]
		if v.Path == "" {
			return errors.New("with invalid view")
		}
		if v.Method == "" {
			v.Method = http.MethodPost
		}
		if v.Prepare == nil {
			// v.Prepare = DefaultPrepareQuery
		}
		// r.Handle(v.Method, filepath.Join(p, v.Path), func(ctx *fiber.Ctx) {
		// 	handleQueryObject(ctx, obj, v.Prepare)
		// })
	}

	return nil
}

func (obj *WebObject) BuildPrimaryPath(prefix string) string {
	var primaryKeyPath []string
	for _, v := range obj.uniqueKeys {
		primaryKeyPath = append(primaryKeyPath, ":"+v.JSONName)
	}
	return filepath.Join(prefix, filepath.Join(primaryKeyPath...))
}

func (obj *WebObject) getPrimaryValues(c *fiber.Ctx) ([]string, error) {
	var result []string
	for _, field := range obj.uniqueKeys {
		v := c.Params(field.JSONName)
		if v == "" {
			return nil, fmt.Errorf("invalid primary: %s", field.JSONName)
		}
		result = append(result, v)
	}
	return result, nil
}

func (obj *WebObject) buildPrimaryCondition(db *gorm.DB, keys []string) *gorm.DB {
	var tx *gorm.DB
	for i := 0; i < len(obj.uniqueKeys); i++ {
		colName := obj.uniqueKeys[i].Name
		col := db.NamingStrategy.ColumnName(obj.tableName, colName)
		tx = db.Where(col, keys[i])
	}
	return tx
}

/*
Check Go type corresponds to JSON type.
- float64, for JSON numbers
- string, for JSON strings
- []any, for JSON arrays
- map[string]any, for JSON objects
- nil, for JSON null
*/
func (obj *WebObject) checkType(db *gorm.DB, key string, value any) (string, bool, error) {
	targetKind, ok := obj.jsonToKinds[key]
	if !ok {
		return "", false, nil
	}

	fieldName, ok := obj.jsonToFields[key]
	if !ok {
		return "", false, nil
	}

	valueKind := reflect.TypeOf(value).Kind()
	var result bool

	switch targetKind {
	case reflect.Struct, reflect.Slice: // time.Time, associated structures
		result = true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		result = valueKind == reflect.Float64
	default:
		result = targetKind == valueKind
	}

	fieldName = db.NamingStrategy.ColumnName(obj.tableName, fieldName)
	if !result {
		return fieldName, false, fmt.Errorf("%s type not match", key)
	}
	return fieldName, true, nil
}

// Build fill the properties of obj.
func (obj *WebObject) Build() error {
	rt := reflect.TypeOf(obj.Model)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	obj.modelElem = rt
	obj.tableName = obj.modelElem.Name()

	if obj.Name == "" {
		obj.Name = strings.ToLower(obj.tableName)
	}

	obj.jsonToFields = make(map[string]string)
	obj.jsonToKinds = make(map[string]reflect.Kind)
	obj.parseFields(obj.modelElem)

	if obj.primaryKeys != nil {
		obj.uniqueKeys = obj.primaryKeys
	}

	if len(obj.uniqueKeys) <= 0 && len(obj.primaryKeys) <= 0 {
		return fmt.Errorf("%s not has primaryKey", obj.Name)
	}
	return nil
}

// parseFields parse the following properties according to struct tag:
// - jsonToFields, jsonToKinds, primaryKeyName, primaryKeyJsonName
func (obj *WebObject) parseFields(rt reflect.Type) {
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)

		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			obj.parseFields(f.Type)
			continue
		}

		jsonTag := strings.TrimSpace(strings.Split(f.Tag.Get("json"), ",")[0])
		if jsonTag == "" {
			obj.jsonToFields[f.Name] = f.Name

			kind := f.Type.Kind()
			if kind == reflect.Ptr {
				kind = f.Type.Elem().Kind()
			}
			obj.jsonToKinds[f.Name] = kind
		} else if jsonTag != "-" {
			obj.jsonToFields[jsonTag] = f.Name

			kind := f.Type.Kind()
			if kind == reflect.Ptr {
				kind = f.Type.Elem().Kind()
			}
			obj.jsonToKinds[jsonTag] = kind
		}

		gormTag := strings.ToLower(f.Tag.Get("gorm"))
		if gormTag == "-" {
			continue
		}
		pkField := WebObjectPrimaryField{
			Name:      f.Name,
			JSONName:  strings.Split(jsonTag, ",")[0],
			Kind:      f.Type.Kind(),
			IsPrimary: strings.Contains(gormTag, "primarykey"),
		}

		if pkField.JSONName == "" {
			pkField.JSONName = pkField.Name
		}

		if pkField.IsPrimary {
			obj.primaryKeys = append(obj.primaryKeys, pkField)
		} else if strings.Contains(gormTag, "unique") {
			obj.uniqueKeys = append(obj.uniqueKeys, pkField)
		}
	}
}
