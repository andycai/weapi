package entity

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/andycai/weapi/administrator/components/entity"
	"github.com/andycai/weapi/object"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterObjects(r fiber.Router, objs []object.WebObject) {
	for idx := range objs {
		obj := &objs[idx]
		err := RegisterObject(obj, r)
		if err != nil {
			// log.Fatalf("RegisterObject [%s] fail %v\n", obj.Name, err)
		}
	}
}

func RegisterObject(obj *object.WebObject, r fiber.Router) error {
	if err := Build(obj); err != nil {
		return err
	}

	p := obj.Name
	allowMethods := obj.AllowMethods
	if allowMethods == 0 {
		allowMethods = object.GET | object.CREATE | object.EDIT | object.DELETE | object.QUERY
	}

	primaryKeyPath := BuildPrimaryPath(obj, p)
	if allowMethods&object.GET != 0 {
		r.Get(primaryKeyPath, func(c *fiber.Ctx) error {
			handleGetObject(c, obj)
			return nil
		})
	}
	if allowMethods&object.CREATE != 0 {
		r.Put(p, func(c *fiber.Ctx) error {
			handleCreateObject(c, obj)
			return nil
		})
	}
	if allowMethods&object.EDIT != 0 {
		r.Patch(primaryKeyPath, func(c *fiber.Ctx) error {
			handleEditObject(c, obj)
			return nil
		})
	}

	if allowMethods&object.DELETE != 0 {
		r.Delete(primaryKeyPath, func(c *fiber.Ctx) error {
			handleDeleteObject(c, obj)
			return nil
		})
	}

	if allowMethods&object.QUERY != 0 {
		r.Post(p, func(c *fiber.Ctx) error {
			handleQueryObject(c, obj, entity.DefaultPrepareQuery)
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
			v.Prepare = entity.DefaultPrepareQuery
		}
		r.Add(v.Method, filepath.Join(p, v.Path), func(ctx *fiber.Ctx) error {
			handleQueryObject(ctx, obj, v.Prepare)
			return nil
		})
	}

	return nil
}

func BuildPrimaryPath(obj *object.WebObject, prefix string) string {
	var primaryKeyPath []string
	for _, v := range obj.UniqueKeys {
		primaryKeyPath = append(primaryKeyPath, ":"+v.JSONName)
	}
	return filepath.Join(prefix, filepath.Join(primaryKeyPath...))
}

func getPrimaryValues(obj *object.WebObject, c *fiber.Ctx) ([]string, error) {
	var result []string
	for _, field := range obj.UniqueKeys {
		v := c.Params(field.JSONName)
		if v == "" {
			return nil, fmt.Errorf("invalid primary: %s", field.JSONName)
		}
		result = append(result, v)
	}
	return result, nil
}

func buildPrimaryCondition(obj *object.WebObject, keys []string) *gorm.DB {
	var tx *gorm.DB
	for i := 0; i < len(obj.UniqueKeys); i++ {
		colName := obj.UniqueKeys[i].Name
		col := db.NamingStrategy.ColumnName(obj.TableName, colName)
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
func checkType(obj *object.WebObject, key string, value any) (string, bool, error) {
	targetKind, ok := obj.JsonToKinds[key]
	if !ok {
		return "", false, nil
	}

	fieldName, ok := obj.JsonToFields[key]
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

	fieldName = db.NamingStrategy.ColumnName(obj.TableName, fieldName)
	if !result {
		return fieldName, false, fmt.Errorf("%s type not match", key)
	}
	return fieldName, true, nil
}

// Build fill the properties of obj.
func Build(obj *object.WebObject) error {
	rt := reflect.TypeOf(obj.Model)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	obj.ModelElem = rt
	obj.TableName = obj.ModelElem.Name()

	if obj.Name == "" {
		obj.Name = strings.ToLower(obj.TableName)
	}

	obj.JsonToFields = make(map[string]string)
	obj.JsonToKinds = make(map[string]reflect.Kind)
	parseFields(obj, obj.ModelElem)

	if obj.PrimaryKeys != nil {
		obj.UniqueKeys = obj.PrimaryKeys
	}

	if len(obj.UniqueKeys) <= 0 && len(obj.PrimaryKeys) <= 0 {
		return fmt.Errorf("%s not has primaryKey", obj.Name)
	}
	return nil
}

// parseFields parse the following properties according to struct tag:
// - JsonToFields, JsonToKinds, primaryKeyName, primaryKeyJsonName
func parseFields(obj *object.WebObject, rt reflect.Type) {
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)

		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			parseFields(obj, f.Type)
			continue
		}

		jsonTag := strings.TrimSpace(strings.Split(f.Tag.Get("json"), ",")[0])
		if jsonTag == "" {
			obj.JsonToFields[f.Name] = f.Name

			kind := f.Type.Kind()
			if kind == reflect.Ptr {
				kind = f.Type.Elem().Kind()
			}
			obj.JsonToKinds[f.Name] = kind
		} else if jsonTag != "-" {
			obj.JsonToFields[jsonTag] = f.Name

			kind := f.Type.Kind()
			if kind == reflect.Ptr {
				kind = f.Type.Elem().Kind()
			}
			obj.JsonToKinds[jsonTag] = kind
		}

		gormTag := strings.ToLower(f.Tag.Get("gorm"))
		if gormTag == "-" {
			continue
		}
		pkField := object.WebObjectPrimaryField{
			Name:      f.Name,
			JSONName:  strings.Split(jsonTag, ",")[0],
			Kind:      f.Type.Kind(),
			IsPrimary: strings.Contains(gormTag, "primarykey"),
		}

		if pkField.JSONName == "" {
			pkField.JSONName = pkField.Name
		}

		if pkField.IsPrimary {
			obj.PrimaryKeys = append(obj.PrimaryKeys, pkField)
		} else if strings.Contains(gormTag, "unique") {
			obj.UniqueKeys = append(obj.UniqueKeys, pkField)
		}
	}
}

//#region api handler

func handleGetObject(c *fiber.Ctx, obj *object.WebObject) {
	keys, err := getPrimaryValues(obj, c)
	if err != nil {
		// AbortWithJSONError(c, http.StatusBadRequest, err)
		return
	}
	// db := getDbConnection(c, obj.GetDB, false)
	// the real name of the primaryKey column
	val := reflect.New(obj.ModelElem).Interface()
	result := buildPrimaryCondition(obj, keys).Take(&val)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// AbortWithJSONError(c, http.StatusNotFound, errors.New("not found"))
		} else {
			// AbortWithJSONError(c, http.StatusInternalServerError, result.Error)
		}
		return
	}

	if obj.BeforeRender != nil {
		rr, err := obj.BeforeRender(c, val)
		if err != nil {
			// AbortWithJSONError(c, http.StatusInternalServerError, err)
			return
		}

		// if c.Writer.Written() || c.Writer.Status() != http.StatusOK {
		// if body has written, return
		// return
		// }

		if rr != nil {
			val = rr
		}
	}

	c.JSON(val)
}

func handleCreateObject(c *fiber.Ctx, obj *object.WebObject) {
	val := reflect.New(obj.ModelElem).Interface()

	if c.Request().Header.ContentLength() > 0 {
		if err := c.BodyParser(&val); err != nil {
			// AbortWithJSONError(c, http.StatusBadRequest, err)
			return
		}
	}

	// db := getDbConnection(c, obj.GetDB, true)
	if obj.BeforeCreate != nil {
		if err := obj.BeforeCreate(c, val); err != nil {
			// AbortWithJSONError(c, http.StatusBadRequest, err)
			return
		}
	}

	result := db.Create(val)
	if result.Error != nil {
		// AbortWithJSONError(c, http.StatusInternalServerError, result.Error)
		return
	}

	c.JSON(val)
}

func handleEditObject(c *fiber.Ctx, obj *object.WebObject) {
	keys, err := getPrimaryValues(obj, c)
	if err != nil {
		// AbortWithJSONError(c, http.StatusBadRequest, err)
		return
	}

	var inputVals map[string]any
	if err := c.BodyParser(&inputVals); err != nil {
		// AbortWithJSONError(c, http.StatusBadRequest, err)
		return
	}

	// db := getDbConnection(c, obj.GetDB, false)

	var vals map[string]any = map[string]any{}

	// can't edit primaryKey
	for _, k := range obj.UniqueKeys {
		delete(inputVals, k.JSONName)
	}

	for k, v := range inputVals {
		if v == nil {
			continue
		}

		fieldName, ok, err := checkType(obj, k, v)
		if err != nil {
			// AbortWithJSONError(c, http.StatusBadRequest, fmt.Errorf("%s type not match", k))
			return
		}
		if !ok { // ignore invalid field
			continue
		}
		vals[fieldName] = v
	}

	if len(obj.Editables) > 0 {
		stripVals := make(map[string]any)
		for _, k := range obj.Editables {
			k = db.NamingStrategy.ColumnName(obj.TableName, k)
			if v, ok := vals[k]; ok {
				stripVals[k] = v
			}
		}
		vals = stripVals
	} else {
		vals = map[string]any{}
	}

	if len(vals) == 0 {
		// AbortWithJSONError(c, http.StatusBadRequest, errors.New("not changed"))
		return
	}
	// db = buildPrimaryCondition(obj, db.Table(db.NamingStrategy.TableName(obj.TableName)), keys)
	db = buildPrimaryCondition(obj, keys)

	if obj.BeforeUpdate != nil {
		val := reflect.New(obj.ModelElem).Interface()
		tx := db.Session(&gorm.Session{})
		if err := tx.First(val).Error; err != nil {
			// AbortWithJSONError(c, http.StatusNotFound, errors.New("not found"))
			return
		}
		if err := obj.BeforeUpdate(c, val, inputVals); err != nil {
			// AbortWithJSONError(c, http.StatusBadRequest, err)
			return
		}
	}

	result := db.Updates(vals)
	if result.Error != nil {
		// AbortWithJSONError(c, http.StatusInternalServerError, result.Error)
		return
	}

	c.JSON(true)
}

func handleDeleteObject(c *fiber.Ctx, obj *object.WebObject) {
	keys, err := getPrimaryValues(obj, c)
	if err != nil {
		// AbortWithJSONError(c, http.StatusBadRequest, err)
		return
	}

	// db := getDbConnection(c, obj.GetDB, false)
	val := reflect.New(obj.ModelElem).Interface()

	r := buildPrimaryCondition(obj, keys).Session(&gorm.Session{}).First(val)

	// for gorm delete hook, need to load model first.
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			// AbortWithJSONError(c, http.StatusNotFound, errors.New("not found"))
		} else {
			// AbortWithJSONError(c, http.StatusInternalServerError, r.Error)
		}
		return
	}

	if obj.BeforeDelete != nil {
		if err := obj.BeforeDelete(c, val); err != nil {
			// AbortWithJSONError(c, http.StatusBadRequest, err)
			return
		}
	}

	r = db.Delete(val)
	if r.Error != nil {
		// AbortWithJSONError(c, http.StatusInternalServerError, r.Error)
		return
	}

	c.JSON(true)
}

func handleQueryObject(c *fiber.Ctx, obj *object.WebObject, prepareQuery object.PrepareQuery) {
	db, form, err := prepareQuery(db, c)
	if err != nil {
		// AbortWithJSONError(c, http.StatusBadRequest, err)
		return
	}

	namer := db.NamingStrategy

	// Use struct{} makes map like set.
	var filterFields = make(map[string]struct{})
	for _, k := range obj.Filterables {
		filterFields[k] = struct{}{}
	}

	if len(filterFields) > 0 {
		var stripFilters []object.Filter
		for i := 0; i < len(form.Filters); i++ {
			filter := form.Filters[i]
			// Struct must has this field.
			field, ok := obj.JsonToFields[filter.Name]
			if !ok {
				continue
			}
			if _, ok := filterFields[field]; !ok {
				continue
			}

			if f, ok := obj.ModelElem.FieldByName(field); ok {
				var typeName string = f.Type.Name()
				if f.Type.Kind() == reflect.Ptr {
					typeName = f.Type.Elem().Name()
				}
				filter.IsTimeType = typeName == "Time" || typeName == "NullTime" || typeName == "DeletedAt"
			}
			filter.Name = namer.ColumnName(obj.TableName, field)
			stripFilters = append(stripFilters, filter)
		}
		form.Filters = stripFilters
	} else {
		form.Filters = []object.Filter{}
	}

	var orderFields = make(map[string]struct{})
	for _, k := range obj.Orderables {
		orderFields[k] = struct{}{}
	}
	if len(orderFields) > 0 {
		var stripOrders []object.Order
		for i := 0; i < len(form.Orders); i++ {
			order := form.Orders[i]
			field, ok := obj.JsonToFields[order.Name]
			if !ok {
				continue
			}
			if _, ok := orderFields[field]; !ok {
				continue
			}
			order.Name = namer.ColumnName(obj.TableName, order.Name)
			stripOrders = append(stripOrders, order)
		}
		form.Orders = stripOrders
	} else {
		form.Orders = []object.Order{}
	}

	if form.Keyword != "" {
		form.SearchFields = []string{}
		for _, v := range obj.Searchables {
			form.SearchFields = append(form.SearchFields, namer.ColumnName(obj.TableName, v))
		}
	}

	if len(form.ViewFields) > 0 {
		var stripViewFields []string
		for _, v := range form.ViewFields {
			stripViewFields = append(stripViewFields, namer.ColumnName(obj.TableName, v))
		}
		form.ViewFields = stripViewFields
	}

	r, err := queryObjects(obj, c, form)
	if err != nil {
		// AbortWithJSONError(c, http.StatusBadRequest, err)
		return
	}

	if obj.BeforeQueryRender != nil {
		obj, err := obj.BeforeQueryRender(c, &r)
		if err != nil {
			// AbortWithJSONError(c, http.StatusBadRequest, err)
			return
		}

		// if c.Writer.Written() || c.Writer.Status() != http.StatusOK {
		// if body has written, return
		// return
		// }

		if obj != nil {
			c.JSON(obj)
			return
		}
	}
	c.JSON(r)
}

func castTime(value any) any {
	if tv, ok := value.(string); ok {
		for _, tf := range []string{time.RFC3339, time.RFC3339Nano, "2006-01-02 15:04:05", "2006-01-02", time.RFC1123} {
			t, err := time.Parse(tf, tv)
			if err == nil {
				return t
			}
		}
	}
	return value
}

func queryObjects(obj *object.WebObject, ctx *fiber.Ctx, form *object.QueryForm) (r object.QueryResult, err error) {
	tblName := db.NamingStrategy.TableName(obj.TableName)

	for _, v := range form.Filters {
		if q := v.GetQuery(); q != "" {
			if v.Op == object.FilterOpLike {
				if kws, ok := v.Value.([]any); ok {
					qs := []string{}
					for _, kw := range kws {
						k := fmt.Sprintf("\"%%%s%%\"", strings.ReplaceAll(kw.(string), "\"", "\\\""))
						q := fmt.Sprintf("`%s`.`%s` LIKE %s", tblName, v.Name, k)
						qs = append(qs, q)
					}
					db = db.Where(strings.Join(qs, " OR "))
				} else {
					db = db.Where(fmt.Sprintf("`%s`.%s", tblName, q), fmt.Sprintf("%%%s%%", v.Value))
				}
			} else if v.Op == object.FilterOpBetween {
				vt := reflect.ValueOf(v.Value)
				if vt.Kind() != reflect.Slice && vt.Len() != 2 {
					return r, fmt.Errorf("invalid between value, must be slice with 2 elements")
				}

				leftValue := vt.Index(0).Interface()
				rightValue := vt.Index(1).Interface()
				if v.IsTimeType {
					leftValue = castTime(leftValue)
					rightValue = castTime(rightValue)
				}
				db = db.Where(fmt.Sprintf("`%s`.%s", tblName, q), leftValue, rightValue)
			} else {
				if v.IsTimeType {
					v.Value = castTime(v.Value)
				}
				db = db.Where(fmt.Sprintf("`%s`.%s", tblName, q), v.Value)
			}
		}
	}

	for _, v := range form.Orders {
		if q := v.GetQuery(); q != "" {
			db = db.Order(fmt.Sprintf("%s.%s", tblName, q))
		}
	}

	if form.Keyword != "" && len(form.SearchFields) > 0 {
		var query []string
		for _, v := range form.SearchFields {
			query = append(query, fmt.Sprintf("`%s`.`%s` LIKE @keyword", tblName, v))
		}
		searchKey := strings.Join(query, " OR ")
		db = db.Where(searchKey, sql.Named("keyword", "%"+form.Keyword+"%"))
	}

	if len(form.ViewFields) > 0 {
		db = db.Select(form.ViewFields)
	}

	r.Pos = form.Pos
	r.Limit = form.Limit
	r.Keyword = form.Keyword

	var c int64
	if err := db.Table(tblName).Count(&c).Error; err != nil {
		return r, err
	}
	if c <= 0 {
		return r, nil
	}
	r.TotalCount = int(c)

	vals := reflect.New(reflect.SliceOf(obj.ModelElem))
	result := db.Offset(form.Pos).Limit(form.Limit).Find(vals.Interface())
	if result.Error != nil {
		return r, result.Error
	}

	r.Items = make([]any, 0, vals.Elem().Len())
	for i := 0; i < vals.Elem().Len(); i++ {
		modelObj := vals.Elem().Index(i).Addr().Interface()
		if obj.BeforeRender != nil {
			rr, err := obj.BeforeRender(ctx, modelObj)
			if err != nil {
				return r, err
			}
			if rr != nil {
				// if BeforeRender return not nil, then use it as result
				modelObj = rr
			}
		}
		r.Items = append(r.Items, modelObj)
	}
	r.Pos += int(len(r.Items))
	return r, nil
}

//#endregion
