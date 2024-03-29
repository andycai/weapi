package site

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/andycai/weapi/admin/content"
	"github.com/andycai/weapi/admin/user"
	"github.com/andycai/weapi/constant"
	"github.com/andycai/weapi/log"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/inflection"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetPageContext() map[string]any {
	return map[string]any{
		"siteurl":            user.GetValue(constant.KEY_SITE_URL),
		"sitename":           user.GetValue(constant.KEY_SITE_NAME),
		"copyright":          user.GetValue(constant.KEY_SITE_COPYRIGHT),
		"siteadmin":          user.GetValue(constant.KEY_SITE_ADMIN),
		"keywords":           user.GetValue(constant.KEY_SITE_KEYWORDS),
		"description":        user.GetValue(constant.KEY_SITE_DESCRIPTION),
		"ga":                 user.GetValue(constant.KEY_SITE_GA),
		"logo_url":           user.GetValue(constant.KEY_SITE_LOGO_URL),
		"favicon_url":        user.GetValue(constant.KEY_SITE_FAVICON_URL),
		"terms_url":          user.GetValue(constant.KEY_SITE_TERMS_URL),
		"privacy_url":        user.GetValue(constant.KEY_SITE_PRIVACY_URL),
		"signin_url":         user.GetValue(constant.KEY_SITE_SIGNIN_URL),
		"signup_url":         user.GetValue(constant.KEY_SITE_SIGNUP_URL),
		"logout_url":         user.GetValue(constant.KEY_SITE_LOGOUT_URL),
		"reset_password_url": user.GetValue(constant.KEY_SITE_RESET_PASSWORD_URL),
		"login_next":         user.GetValue(constant.KEY_SITE_LOGIN_NEXT),
		"slogan":             user.GetValue(constant.KEY_SITE_SLOGAN),
		"user_id_type":       user.GetValue(constant.KEY_SITE_USER_ID_TYPE),
		"dashboard":          user.GetValue(constant.KEY_ADMIN_DASHBOARD),
	}
}

func decorateHandler(obj *model.AdminObject, handler func(obj *model.AdminObject, c *fiber.Ctx) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return handler(obj, c)
	}
}

func BuildAdminObjects(r fiber.Router, objs []model.AdminObject) []*model.AdminObject {
	handledObjects := make([]*model.AdminObject, 0)
	exists := make(map[string]bool)
	for idx := range objs {
		obj := &objs[idx]
		err := Build(obj)
		if err != nil {
			log.Infof("Build admin object fail, ignore %s, %s, err: %s", obj.Group, obj.Name, err)
			continue
		}

		if _, ok := exists[obj.Path]; ok {
			log.Infof("Ignore exist admin object %s, %s", obj.Group, obj.Name)
			continue
		}

		objr := r.Group(obj.Path)
		obj.Path = path.Join("", obj.Path) + "/"
		for idx := range obj.Fields {
			f := &obj.Fields[idx]
			if f.Foreign == nil {
				continue
			}
			f.Foreign.Path = path.Join("", f.Foreign.Path) + "/"
		}

		RegisterAdminRouter(obj, objr)
		handledObjects = append(handledObjects, obj)
	}
	return handledObjects
}

func BuildPermissions(obj *model.AdminObject, user *model.User) {
	obj.Permissions = map[string]bool{}
	if user.IsSuperUser {
		obj.Permissions["can_create"] = true
		obj.Permissions["can_update"] = true
		obj.Permissions["can_delete"] = true
		obj.Permissions["can_action"] = true
		return
	}

	//TODO: build permissions with group settings
	obj.Permissions["can_create"] = true
	obj.Permissions["can_update"] = true
	obj.Permissions["can_delete"] = true
	obj.Permissions["can_action"] = true
}

// RegisterAdminRouter registers admin routes
//
//   - POST /admin/{objectslug} -> Query objects
//   - PUT /admin/{objectslug} -> Create One
//   - PATCH /admin/{objectslug}} -> Update One
//   - DELETE /admin/{objectslug} -> Delete One
//   - POST /admin/{objectslug}/:name -> Action
func RegisterAdminRouter(obj *model.AdminObject, r fiber.Router) {
	r = r.Use(func(ctx *fiber.Ctx) error {
		if obj.AccessCheck != nil {
			err := obj.AccessCheck(ctx, obj)
			if err != nil {
				return err
			}
		}
		return ctx.Next()
	})

	r.Post("/", decorateHandler(obj, handleQueryOrGetOne))
	r.Put("/", decorateHandler(obj, handleCreate))
	r.Patch("/", decorateHandler(obj, handleUpdate))
	r.Delete("/", decorateHandler(obj, handleDelete))
	r.Post("/:name", decorateHandler(obj, handleAction))
}

func asColNames(obj *model.AdminObject, db *gorm.DB, fields []string) []string {
	for i := 0; i < len(fields); i++ {
		fields[i] = db.NamingStrategy.ColumnName(obj.TableName, fields[i])
	}
	return fields
}

// Build fill the properties of obj.
func Build(obj *model.AdminObject) error {
	if obj.Path == "" {
		obj.Path = strings.ToLower(obj.Name)
	}

	if obj.Path == "_" || obj.Path == "" {
		return fmt.Errorf("invalid path")
	}

	rt := reflect.TypeOf(obj.Model)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	obj.ModelElem = rt
	obj.TableName = db.NamingStrategy.TableName(rt.Name())
	obj.PluralName = inflection.Plural(obj.Name)
	obj.Shows = asColNames(obj, db, obj.Shows)
	obj.Editables = asColNames(obj, db, obj.Editables)
	obj.Orderables = asColNames(obj, db, obj.Orderables)
	obj.Searchables = asColNames(obj, db, obj.Searchables)
	obj.Filterables = asColNames(obj, db, obj.Filterables)
	obj.Requireds = asColNames(obj, db, obj.Requireds)
	obj.PrimaryKeyMaping = map[string]string{}

	for idx := range obj.Orders {
		o := &obj.Orders[idx]
		o.Name = db.NamingStrategy.ColumnName(obj.TableName, o.Name)
	}

	obj.Ignores = map[string]bool{}
	err := parseFields(obj, db, rt)
	if err != nil {
		return err
	}
	if len(obj.PrimaryKeys) <= 0 && len(obj.UniqueKeys) <= 0 {
		return fmt.Errorf("%s not has primaryKey or uniqueKeys", obj.Name)
	}

	for idx := range obj.Actions {
		action := &obj.Actions[idx]
		if action.Name == "" {
			continue
		}
		if action.Path == "" {
			action.Path = strings.ToLower(action.Name)
		}
	}
	return nil
}

func parseFields(obj *model.AdminObject, db *gorm.DB, rt reflect.Type) error {
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)

		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			parseFields(obj, db, f.Type)
			continue
		}

		if f.Type.Kind() == reflect.Chan || f.Type.Kind() == reflect.Func || !f.IsExported() {
			continue
		}

		gormTag := strings.ToLower(f.Tag.Get("gorm"))
		field := model.AdminField{
			Name:      db.NamingStrategy.ColumnName(obj.TableName, f.Name),
			Tag:       gormTag,
			ElemType:  f.Type,
			FieldName: f.Name,
			Label:     f.Tag.Get("label"),
			NotColumn: gormTag == "-",
		}
		if field.ElemType.Kind() == reflect.Ptr {
			field.ElemType = field.ElemType.Elem()
		}
		if field.Label == "" {
			field.Label = strings.ReplaceAll(field.Name, "_", " ")
		}

		field.Label = cases.Title(language.Und).String(field.Label)

		switch f.Type.Kind() {
		case reflect.Ptr:
			field.Type = f.Type.Elem().Name()
			field.CanNull = true
			field.IsPtr = true
		case reflect.Slice:
			field.Type = f.Type.Elem().Name()
			field.CanNull = true
			field.IsArray = true
		default:
			field.Type = f.Type.Name()
		}

		if strings.Contains(gormTag, "primarykey") {
			field.Primary = true
			//obj.primaryKeys = append(obj.primaryKeys, field.Name)
			if strings.Contains(field.Type, "int") {
				field.IsAutoID = true
			}
		}

		if strings.Contains(gormTag, "primarykey") || strings.Contains(gormTag, "unique") {
			// hint foreignField
			keyName := field.Name
			if strings.HasSuffix(f.Name, "ID") {
				n := f.Name[:len(f.Name)-2]
				if ff, ok := rt.FieldByName(n); ok {
					if ff.Type.Kind() == reflect.Struct || (ff.Type.Kind() == reflect.Ptr && ff.Type.Elem().Kind() == reflect.Struct) {
						keyName = db.NamingStrategy.ColumnName(obj.TableName, ff.Name)
					}
				}
				obj.PrimaryKeyMaping[keyName] = field.Name
			}
			if strings.Contains(gormTag, "primarykey") {
				obj.PrimaryKeys = append(obj.PrimaryKeys, keyName)
			} else {
				obj.UniqueKeys = append(obj.UniqueKeys, keyName)
			}
		}

		foreignKey := ""
		// ignore `belongs to` and `has one` relation
		if f.Type.Kind() == reflect.Struct || (f.Type.Kind() == reflect.Ptr && f.Type.Elem().Kind() == reflect.Struct) {
			hintForeignKey := fmt.Sprintf("%sID", f.Name)
			if _, ok := rt.FieldByName(hintForeignKey); ok {
				foreignKey = hintForeignKey
			}
		}
		if strings.Contains(gormTag, "foreignkey") {
			//extract foreign key from gorm tag with regex
			//example: foreignkey:UserRefer
			var re = regexp.MustCompile(`foreignkey:([a-zA-Z0-9]+)`)
			matches := re.FindStringSubmatch(gormTag)
			if len(matches) > 1 {
				foreignKey = strings.TrimSpace(matches[1])
			}
		}

		if foreignKey != "" {
			obj.Ignores[foreignKey] = true
			for k := range obj.Fields {
				if obj.Fields[k].FieldName == foreignKey {
					// remove foreign key from fields
					obj.Fields = append(obj.Fields[:k], obj.Fields[k+1:]...)
					break
				}
			}

			field.Foreign = &model.AdminForeign{
				Field:      db.NamingStrategy.ColumnName(obj.TableName, foreignKey),
				Path:       strings.ToLower(f.Type.Name()),
				ForeignKey: foreignKey,
				FieldName:  f.Name,
			}
		}

		if field.Type == "DeletedAt" {
			obj.MarkDeletedField = field.Name
		}

		if field.Type == "NullTime" || field.Type == "Time" || field.Type == "DeletedAt" {
			field.Type = "datetime"
		}

		if field.Type == "DeletedAt" || strings.HasPrefix("Null", field.Type) {
			field.CanNull = true
		}

		if attr, ok := obj.Attributes[f.Name]; ok {
			field.Attribute = &attr
		}
		obj.Fields = append(obj.Fields, field)
	}
	return nil
}

func formatAsInt64(v any) int64 {
	srcKind := reflect.ValueOf(v).Kind()
	switch srcKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(v).Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(reflect.ValueOf(v).Uint())
	case reflect.Float32, reflect.Float64:
		return int64(reflect.ValueOf(v).Float())
	case reflect.String:
		if v.(string) == "" {
			return 0
		}
		if i, err := strconv.ParseInt(v.(string), 10, 64); err == nil {
			return i
		}
	}
	return 0
}

func convertValue(ElemType reflect.Type, source any) (any, error) {
	srcType := reflect.TypeOf(source)
	if srcType == ElemType {
		return source, nil
	}

	// if srcType.Kind() == reflect.Array || srcType.Kind() == reflect.Slice || srcType.Kind() == reflect.Map {
	// 	return source, nil
	// }

	var targetType reflect.Type = ElemType
	var err error
	switch ElemType.Name() {
	case "int", "int8", "int16", "int32", "int64":
		v := formatAsInt64(source)
		return reflect.ValueOf(v).Convert(targetType).Interface(), nil
	case "uint", "uint8", "uint16", "uint32", "uint64":
		v := formatAsInt64(source)
		return reflect.ValueOf(v).Convert(targetType).Interface(), nil
	case "float32", "float64":
		v, err := strconv.ParseFloat(fmt.Sprintf("%v", source), 64)
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(v).Convert(targetType).Interface(), nil
	case "bool":
		val := fmt.Sprintf("%v", source)
		if val == "on" {
			val = "true"
		} else if val == "off" {
			val = "false"
		} else if val == "" {
			val = "false"
		}

		v, err := strconv.ParseBool(val)
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(v).Interface(), nil
	case "string":
		return fmt.Sprintf("%v", source), nil
	case "NullTime":
		tv, ok := source.(string)
		if tv == "" || !ok {
			return &sql.NullTime{}, nil
		} else {
			for _, tf := range []string{time.RFC3339, time.RFC3339Nano, "2006-01-02 15:04:05", "2006-01-02", time.RFC1123} {
				t, err := time.Parse(tf, tv)
				if err == nil {
					return &sql.NullTime{Time: t, Valid: true}, nil
				}
			}
		}
		return nil, fmt.Errorf("invalid datetime format %v", source)
	case "Time":
		tv, ok := source.(string)
		if tv == "" || !ok {
			return &time.Time{}, nil
		} else {
			for _, tf := range []string{time.RFC3339, time.RFC3339Nano, "2006-01-02 15:04:05", "2006-01-02", time.RFC1123} {
				t, err := time.Parse(tf, tv)
				if err == nil {
					return &t, nil
				}
			}
		}
		return nil, fmt.Errorf("invalid datetime format %v", source)
	default:
		var data []byte
		data, err = json.Marshal(source)
		if err != nil {
			return nil, err
		}
		value := reflect.New(targetType).Interface()
		err = json.Unmarshal(data, value)
		return value, err
	}
}

func UnmarshalFrom(obj *model.AdminObject, elemObj reflect.Value, keys, vals map[string]any) (any, error) {
	if len(obj.Editables) > 0 {
		editables := make(map[string]bool)
		for _, v := range obj.Editables {
			editables[v] = true
		}
		for k := range vals {
			if _, ok := editables[k]; !ok {
				delete(vals, k)
			}
		}
	}

	for k, v := range keys {
		// if primary key in editables, then ignore it
		if _, ok := vals[k]; !ok {
			vals[k] = v
		}
	}

	for _, field := range obj.Fields {
		val, ok := vals[field.Name]
		if !ok {
			continue
		}

		if val == nil {
			continue
		}
		var target reflect.Value
		var targetValue reflect.Value
		var targetType = field.ElemType
		if field.Foreign != nil {
			target = elemObj.Elem().FieldByName(field.Foreign.ForeignKey)
			targetType = target.Type()
			if valMap, ok := val.(map[string]any); ok {
				if v, ok := valMap["value"]; ok {
					val = v
				}
			}
		} else {
			target = elemObj.Elem().FieldByName(field.FieldName)
		}

		fieldValue, err := convertValue(targetType, val)
		if err != nil {
			return nil, fmt.Errorf("invalid type: %s except: %s actual: %s error:%v", field.Name, field.Type, reflect.TypeOf(val).Name(), err)
		}
		targetValue = reflect.ValueOf(fieldValue)

		if target.Kind() == reflect.Ptr {
			ptrValue := reflect.New(reflect.PointerTo(field.ElemType))
			ptrValue.Elem().Set(targetValue)
			targetValue = ptrValue.Elem()
		} else {
			if targetValue.Kind() == reflect.Ptr {
				targetValue = targetValue.Elem()
			}
		}
		target.Set(targetValue)
	}
	return elemObj.Interface(), nil
}

func MarshalOne(obj *model.AdminObject, val interface{}) (map[string]any, error) {
	var result = make(map[string]any)
	rv := reflect.ValueOf(val)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	for _, field := range obj.Fields {
		var fieldVal any
		if field.Foreign != nil {
			v := model.AdminValue{
				Value: rv.FieldByName(field.Foreign.ForeignKey).Interface(),
			}
			fv := rv.FieldByName(field.Foreign.FieldName)
			if fv.IsValid() {
				if sv, ok := fv.Interface().(fmt.Stringer); ok {
					v.Label = sv.String()
				}
			}
			if v.Label == "" {
				v.Label = fmt.Sprintf("%v", v.Value)
			}
			fieldVal = v
		} else {
			v := rv.FieldByName(field.FieldName)
			if v.IsValid() {
				fieldVal = v.Interface()
			}
		}
		result[field.Name] = fieldVal
	}
	return result, nil
}

func getPrimaryValues(obj *model.AdminObject, c *fiber.Ctx) map[string]any {
	var result = make(map[string]any)
	hasPrimaryQuery := false
	for _, field := range obj.PrimaryKeys {
		if v := c.Query(field); v != "" {
			result[field] = v
			hasPrimaryQuery = true
		}
	}

	if hasPrimaryQuery {
		return result
	}

	for _, field := range obj.UniqueKeys {
		if key, ok := obj.PrimaryKeyMaping[field]; ok {
			field = key
		}
		if v := c.Query(field); v != "" {
			result[field] = v
		}
	}
	return result
}

func QueryObjects(obj *model.AdminObject, session *gorm.DB, form *model.QueryForm, ctx *fiber.Ctx) (r model.AdminQueryResult, err error) {
	for _, v := range form.Filters {
		if q := v.GetQuery(); q != "" {
			if v.Op == model.FilterOpLike {
				if kws, ok := v.Value.([]any); ok {
					qs := []string{}
					for _, kw := range kws {
						k := fmt.Sprintf("\"%%%s%%\"", strings.ReplaceAll(kw.(string), "\"", "\\\""))
						q := fmt.Sprintf("`%s`.`%s` LIKE %s", obj.TableName, v.Name, k)
						qs = append(qs, q)
					}
					session = session.Where(strings.Join(qs, " OR "))
				} else {
					session = session.Where(fmt.Sprintf("`%s`.%s", obj.TableName, q), fmt.Sprintf("%%%s%%", v.Value))
				}
			} else if v.Op == model.FilterOpBetween {
				vt := reflect.ValueOf(v.Value)
				if vt.Kind() != reflect.Slice && vt.Len() != 2 {
					return r, fmt.Errorf("invalid between value, must be slice with 2 elements")
				}
				session = session.Where(fmt.Sprintf("`%s`.%s", obj.TableName, q), vt.Index(0).Interface(), vt.Index(1).Interface())
			} else {
				session = session.Where(fmt.Sprintf("`%s`.%s", obj.TableName, q), v.Value)
			}
		}
	}

	var orders []model.Order
	if len(form.Orders) > 0 {
		orders = form.Orders
	} else {
		orders = obj.Orders
	}

	for _, v := range orders {
		if q := v.GetQuery(); q != "" && v.Op != "" {
			session = session.Order(fmt.Sprintf("`%s`.%s", obj.TableName, q))
		}
	}

	if form.Keyword != "" && len(obj.Searchables) > 0 {
		var query []string
		for _, v := range obj.Searchables {
			query = append(query, fmt.Sprintf("`%s`.`%s` LIKE @keyword", obj.TableName, v))
		}
		searchKey := strings.Join(query, " OR ")
		session = session.Where(searchKey, sql.Named("keyword", "%"+form.Keyword+"%"))
	}

	r.Pos = form.Pos
	r.Limit = form.Limit
	r.Keyword = form.Keyword

	if obj.MarkDeletedField != "" {
		session = session.Where(fmt.Sprintf("`%s`.`%s` IS NULL", obj.TableName, obj.MarkDeletedField))
	}

	session = session.Table(obj.TableName)

	var c int64
	if err := session.Count(&c).Error; err != nil {
		return r, err
	}
	if c <= 0 {
		return r, nil
	}
	r.TotalCount = int(c)

	selected := []string{}
	for _, v := range obj.Fields {
		if v.NotColumn {
			continue
		}
		if v.Foreign != nil {
			selected = append(selected, v.Foreign.Field)
		} else {
			selected = append(selected, v.Name)
		}
	}

	vals := reflect.New(reflect.SliceOf(obj.ModelElem))
	tx := session.Preload(clause.Associations).Select(selected).Offset(form.Pos)
	if form.Limit > 0 {
		tx = tx.Limit(form.Limit)
	}
	result := tx.Find(vals.Interface())
	if result.Error != nil {
		return r, result.Error
	}

	for i := 0; i < vals.Elem().Len(); i++ {
		modelObj := vals.Elem().Index(i).Addr().Interface()
		r.Objects = append(r.Objects, modelObj)
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
		item, err := MarshalOne(obj, modelObj)
		if err != nil {
			return r, err
		}
		r.Items = append(r.Items, item)
	}
	return r, nil
}

func checkRequired(required []string, inputVals map[string]any) error {
	for _, v := range required {
		if _, ok := inputVals[v]; ok {
			ret := false
			switch inputVals[v].(type) {
			case string:
				if inputVals[v] == "" {
					ret = true
				}
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
				if inputVals[v] == 0 {
					ret = true
				}
			}
			if ret {
				return errors.New(fmt.Sprintf("required field %s", v))
			}
		}
	}
	return nil
}

// DefaultPrepareQuery return default QueryForm.
func DefaultPrepareQuery(c *fiber.Ctx) (*model.QueryForm, error) {
	var form model.QueryForm
	if c.Request().Header.ContentLength() > 0 {
		if err := c.BodyParser(&form); err != nil {
			return nil, err
		}
	}

	if form.Pos < 0 {
		form.Pos = 0
	}
	if form.Limit <= 0 || form.Limit > model.DefaultQueryLimit {
		form.Limit = model.DefaultQueryLimit
	}

	return &form, nil
}

func GetSummary() (result model.SummaryResult) {
	db.Model(&model.Site{}).Count(&result.SiteCount)
	db.Model(&model.Page{}).Count(&result.PageCount)
	db.Model(&model.Post{}).Count(&result.PostCount)
	db.Model(&model.Category{}).Count(&result.CategoryCount)
	db.Model(&model.Media{}).Where("directory", false).Count(&result.MediaCount)

	var latestPosts []model.Post
	db.Order("updated_at desc").Limit(20).Find(&latestPosts)

	for idx := range latestPosts {
		item := content.NewRenderContentFromPost(&latestPosts[idx], false)
		item.PostBody = ""
		result.LatestPosts = append(result.LatestPosts, item)
	}
	return result
}
