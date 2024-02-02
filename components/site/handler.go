package site

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func handleGetObject(c *fiber.Ctx, obj *model.WebObject) error {
	keys, err := getPrimaryValues(obj, c)
	if err != nil {
		return core.Error(c, http.StatusBadRequest, err)
	}
	// the real name of the primaryKey column
	val := reflect.New(obj.ModelElem).Interface()
	result := buildPrimaryCondition(obj, keys).Take(&val)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return core.Error(c, http.StatusNotFound, errors.New("not found"))
		} else {
			return core.Error(c, http.StatusInternalServerError, result.Error)
		}
	}

	if obj.BeforeRender != nil {
		rr, err := obj.BeforeRender(c, val)
		if err != nil {
			return core.Error(c, http.StatusInternalServerError, err)
		}

		// if c.Writer.Written() || c.Writer.Status() != http.StatusOK {
		// if body has written, return
		// return
		// }

		if rr != nil {
			val = rr
		}
	}

	return c.JSON(val)
}

func handleCreateObject(c *fiber.Ctx, obj *model.WebObject) error {
	val := reflect.New(obj.ModelElem).Interface()

	if c.Request().Header.ContentLength() > 0 {
		if err := c.BodyParser(&val); err != nil {
			return core.Error(c, http.StatusBadRequest, err)
		}
	}

	if obj.BeforeCreate != nil {
		if err := obj.BeforeCreate(c, val); err != nil {
			return core.Error(c, http.StatusBadRequest, err)
		}
	}

	result := db.Create(val)
	if result.Error != nil {
		return core.Error(c, http.StatusInternalServerError, result.Error)
	}

	return c.JSON(val)
}

func handleEditObject(c *fiber.Ctx, obj *model.WebObject) error {
	keys, err := getPrimaryValues(obj, c)
	if err != nil {
		return core.Error(c, http.StatusBadRequest, err)
	}

	var inputVals map[string]any
	if err := c.BodyParser(&inputVals); err != nil {
		return core.Error(c, http.StatusBadRequest, err)
	}

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
			return core.Error(c, http.StatusBadRequest, fmt.Errorf("%s type not match", k))
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
		return core.Error(c, http.StatusBadRequest, errors.New("not changed"))
	}
	// db = buildPrimaryCondition(obj, db.Table(db.NamingStrategy.TableName(obj.TableName)), keys)
	db = buildPrimaryCondition(obj, keys)

	if obj.BeforeUpdate != nil {
		val := reflect.New(obj.ModelElem).Interface()
		tx := db.Session(&gorm.Session{})
		if err := tx.First(val).Error; err != nil {
			return core.Error(c, http.StatusNotFound, errors.New("not found"))
		}
		if err := obj.BeforeUpdate(c, val, inputVals); err != nil {
			return core.Error(c, http.StatusBadRequest, err)
		}
	}

	result := db.Updates(vals)
	if result.Error != nil {
		return core.Error(c, http.StatusInternalServerError, result.Error)
	}

	return c.JSON(true)
}

func handleDeleteObject(c *fiber.Ctx, obj *model.WebObject) error {
	keys, err := getPrimaryValues(obj, c)
	if err != nil {
		return core.Error(c, http.StatusBadRequest, err)
	}

	val := reflect.New(obj.ModelElem).Interface()

	r := buildPrimaryCondition(obj, keys).Session(&gorm.Session{}).First(val)

	// for gorm delete hook, need to load model first.
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return core.Error(c, http.StatusNotFound, errors.New("not found"))
		} else {
			return core.Error(c, http.StatusInternalServerError, r.Error)
		}
	}

	if obj.BeforeDelete != nil {
		if err := obj.BeforeDelete(c, val); err != nil {
			return core.Error(c, http.StatusBadRequest, err)
		}
	}

	r = db.Delete(val)
	if r.Error != nil {
		return core.Error(c, http.StatusInternalServerError, r.Error)
	}

	return c.JSON(true)
}

func handleQueryObject(c *fiber.Ctx, obj *model.WebObject, prepareQuery model.PrepareQuery) error {
	form, err := prepareQuery(c)
	if err != nil {
		return core.Error(c, http.StatusBadRequest, err)
	}

	namer := db.NamingStrategy

	// Use struct{} makes map like set.
	var filterFields = make(map[string]struct{})
	for _, k := range obj.Filterables {
		filterFields[k] = struct{}{}
	}

	if len(filterFields) > 0 {
		var stripFilters []model.Filter
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
		form.Filters = []model.Filter{}
	}

	var orderFields = make(map[string]struct{})
	for _, k := range obj.Orderables {
		orderFields[k] = struct{}{}
	}
	if len(orderFields) > 0 {
		var stripOrders []model.Order
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
		form.Orders = []model.Order{}
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
		return core.Error(c, http.StatusBadRequest, err)
	}

	if obj.BeforeQueryRender != nil {
		obj, err := obj.BeforeQueryRender(c, &r)
		if err != nil {
			return core.Error(c, http.StatusBadRequest, err)
		}

		// if c.Writer.Written() || c.Writer.Status() != http.StatusOK {
		// if body has written, return
		// return
		// }

		if obj != nil {
			return c.JSON(obj)
		}
	}
	return c.JSON(r)
}
