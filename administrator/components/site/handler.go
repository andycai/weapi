package site

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/andycai/weapi/administrator/components/content"
	"github.com/andycai/weapi/administrator/components/user"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func handleDashboard(c *fiber.Ctx) error {
	return c.SendFile("templates/admin/index.html")
}

func handleJson(c *fiber.Ctx, adminObjects []*model.AdminObject) error {
	for _, obj := range adminObjects {
		BuildPermissions(obj, user.Current(c))
	}
	return c.JSON(fiber.Map{
		"objects": adminObjects,
		"user":    user.Current(c),
		"site":    GetPageContext(),
	})
}

func HandleAdminIndex(c *fiber.Ctx, objects []*model.AdminObject, buildContext model.AdminBuildContext) {
	var viewObjects []model.AdminObject
	for _, obj := range objects {
		if obj.AccessCheck != nil {
			err := obj.AccessCheck(c, obj)
			if err != nil {
				continue
			}
		}
		val := *obj
		BuildPermissions(obj, user.Current(c))
		viewObjects = append(viewObjects, val)
	}

	siteCtx := GetPageContext()
	if buildContext != nil {
		siteCtx = buildContext(c, siteCtx)
	}

	c.JSON(fiber.Map{
		"objects": viewObjects,
		"user":    user.Current(c),
		"site":    siteCtx,
	})
}

func handleGetOne(obj *model.AdminObject, c *fiber.Ctx) error {
	modelObj := reflect.New(obj.ModelElem).Interface()
	keys := getPrimaryValues(obj, c)
	if len(keys) <= 0 {
		return core.Error(c, http.StatusBadRequest, errors.New("invalid primary key"))
	}

	result := db.Preload(clause.Associations).Where(keys).First(modelObj)

	if result.Error != nil {
		return core.Error(c, http.StatusInternalServerError, result.Error)
	}

	if obj.BeforeRender != nil {
		rr, err := obj.BeforeRender(c, modelObj)
		if err != nil {
			return core.Error(c, http.StatusInternalServerError, err)
		}
		if rr != nil {
			// if BeforeRender return not nil, then use it as result
			modelObj = rr
		}
	}

	data, err := MarshalOne(obj, modelObj)
	if err != nil {
		return core.Error(c, http.StatusInternalServerError, err)
	}

	return c.JSON(data)
}

// Query many objects with filter/limit/offset/order/search
func handleQueryOrGetOne(obj *model.AdminObject, c *fiber.Ctx) error {
	if c.Request().Header.ContentLength() <= 0 {
		return handleGetOne(obj, c)
	}

	form, err := DefaultPrepareQuery(c)
	if err != nil {
		return core.Error(c, http.StatusBadRequest, err)
	}

	if form.ForeignMode {
		form.Limit = 0 // TODO: support foreign mode limit
	}

	r, err := QueryObjects(obj, db, form, c)

	if err != nil {
		return core.Error(c, http.StatusInternalServerError, err)
	}
	if form.ForeignMode {
		var items []map[string]any
		for i := 0; i < len(r.Items); i++ {
			item := map[string]any{}
			var valueVal any
			for _, v := range obj.Fields {
				if v.Primary {
					valueVal = r.Items[i][v.Name]
				}
			}
			if valueVal == nil {
				continue
			}
			item["value"] = valueVal
			iv := r.Objects[i]
			if sv, ok := iv.(fmt.Stringer); ok {
				item["label"] = sv.String()
			} else {
				item["label"] = fmt.Sprintf("%v", valueVal)
			}
			items = append(items, item)
		}
		r.Items = items
	}

	return c.JSON(r)
}

func handleCreate(obj *model.AdminObject, c *fiber.Ctx) error {
	keys := getPrimaryValues(obj, c)
	var vals map[string]any
	if err := c.BodyParser(&vals); err != nil {
		return core.Error(c, http.StatusBadRequest, err)
	}

	err := checkRequired(obj.Requireds, vals)
	if err != nil {
		return core.Error(c, http.StatusBadRequest, err)
	}

	elmObj := reflect.New(obj.ModelElem)
	elm, err := UnmarshalFrom(obj, elmObj, keys, vals)
	if err != nil {
		return core.Error(c, http.StatusBadRequest, err)
	}
	if obj.BeforeCreate != nil {
		if err := obj.BeforeCreate(c, elm); err != nil {
			return core.Error(c, http.StatusBadRequest, err)
		}
	}

	result := db.Create(elm)
	if result.Error != nil {
		return core.Error(c, http.StatusInternalServerError, result.Error)
	}
	if obj.BeforeRender != nil {
		rr, err := obj.BeforeRender(c, elm)
		if err != nil {
			return core.Error(c, http.StatusInternalServerError, err)
		}
		if rr != nil {
			// if BeforeRender return not nil, then use it as result
			elm = rr
		}
	}

	return c.JSON(elm)
}

func handleUpdate(obj *model.AdminObject, c *fiber.Ctx) error {
	keys := getPrimaryValues(obj, c)
	if len(keys) <= 0 {
		return core.Error(c, http.StatusBadRequest, errors.New("invalid primary key"))
	}

	var inputVals map[string]any
	if err := c.BodyParser(&inputVals); err != nil {
		return core.Error(c, http.StatusBadRequest, err)
	}

	err := checkRequired(obj.Requireds, inputVals)
	if err != nil {
		return core.Error(c, http.StatusBadRequest, err)
	}

	elmObj := reflect.New(obj.ModelElem)
	err = db.Where(keys).First(elmObj.Interface()).Error
	if err != nil {
		return core.Error(c, http.StatusNotFound, errors.New("not found"))
	}

	val, err := UnmarshalFrom(obj, elmObj, keys, inputVals)
	if err != nil {
		return core.Error(c, http.StatusBadRequest, err)
	}

	if obj.BeforeUpdate != nil {
		if err := obj.BeforeUpdate(c, val, inputVals); err != nil {
			return core.Error(c, http.StatusBadRequest, err)
		}
	}

	conflictKeys := []clause.Column{}
	if len(obj.PrimaryKeys) > 0 {
		for _, k := range obj.PrimaryKeys {
			conflictKeys = append(conflictKeys, clause.Column{Name: k})
		}
	} else {
		for _, k := range obj.UniqueKeys {
			conflictKeys = append(conflictKeys, clause.Column{Name: k})
		}
	}

	for idx := range conflictKeys {
		k := &conflictKeys[idx]
		if v, ok := obj.PrimaryKeyMaping[k.Name]; ok {
			k.Name = v
		}
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   conflictKeys,
		UpdateAll: true,
	}).Where(keys).Create(val)

	if result.Error != nil {
		return core.Error(c, http.StatusInternalServerError, result.Error)
	}

	return c.JSON(true)
}

func handleDelete(obj *model.AdminObject, c *fiber.Ctx) error {
	keys := getPrimaryValues(obj, c)
	if len(keys) <= 0 {
		return core.Error(c, http.StatusBadRequest, errors.New("invalid primary key"))
	}
	val := reflect.New(obj.ModelElem).Interface()
	r := db.Where(keys).Take(val)

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

	r = db.Where(keys).Delete(val)
	if r.Error != nil {
		return core.Error(c, http.StatusInternalServerError, r.Error)
	}

	return c.JSON(true)
}

func handleAction(obj *model.AdminObject, c *fiber.Ctx) error {
	for _, action := range obj.Actions {
		if action.Path != c.Params("name") {
			continue
		}

		if action.WithoutObject {
			r, err := action.Handler(c, nil)
			if err != nil {
				return core.Error(c, http.StatusInternalServerError, err)
			}
			return c.JSON(r)
		}

		keys := getPrimaryValues(obj, c)
		if len(keys) <= 0 {
			return core.Error(c, http.StatusBadRequest, errors.New("invalid primary key"))
		}
		modelObj := reflect.New(obj.ModelElem).Interface()
		result := db.Where(keys).First(modelObj)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return core.Error(c, http.StatusNotFound, errors.New("not found"))
			} else {
				return core.Error(c, http.StatusInternalServerError, result.Error)
			}
		}
		r, err := action.Handler(c, modelObj)
		if err != nil {
			return core.Error(c, http.StatusInternalServerError, err)
		}
		return c.JSON(r)
	}
	return core.Error(c, http.StatusBadRequest, errors.New("invalid action"))
}

func HandleAdminSummary(c *fiber.Ctx) error {
	result := GetSummary()
	result.CanExport = user.Current(c).IsSuperUser
	return c.JSON(result)
}

func handleGetTags(c *fiber.Ctx) error {
	contentType := c.Params("content_type")
	var form model.TagsForm
	if err := c.BodyParser(&form); err != nil {
		return core.Error(c, http.StatusBadRequest, err)
	}

	tags, err := content.GetTagsByCategory(contentType, &form)
	if err != nil {
		return core.Error(c, http.StatusInternalServerError, err)
	}

	return c.JSON(tags)
}
