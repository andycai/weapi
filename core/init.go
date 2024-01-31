package core

import (
	"sort"

	"github.com/andycai/weapi/object"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var dbMap = map[string]func([]*gorm.DB){}

var routerRootNoCheckMap = map[string]func(fiber.Router){}
var routerRootCheckMap = map[string]func(fiber.Router){}
var routerAPINoCheckMap = map[string]func(fiber.Router){}
var routerAPICheckMap = map[string]func(fiber.Router){}
var routerAdminCheckMap = map[string]func(fiber.Router){}
var adminObjects = []object.AdminObject{}

func RegisterDatabase(dbType string, f func([]*gorm.DB)) {
	if _, ok := dbMap[dbType]; ok {
		panic("duplicate db type: " + dbType)
	}
	dbMap[dbType] = f
}

func RegisterRootNoCheckRouter(routerType string, f func(fiber.Router)) {
	if _, ok := routerRootNoCheckMap[routerType]; ok {
		panic("duplicate router type: " + routerType)
	}
	routerRootNoCheckMap[routerType] = f
}

func RegisterRootCheckRouter(routerType string, f func(fiber.Router)) {
	if _, ok := routerRootCheckMap[routerType]; ok {
		panic("duplicate router type: " + routerType)
	}
	routerRootCheckMap[routerType] = f
}

func RegisterAPINoCheckRouter(routerType string, f func(fiber.Router)) {
	if _, ok := routerAPINoCheckMap[routerType]; ok {
		panic("duplicate router type: " + routerType)
	}
	routerAPINoCheckMap[routerType] = f
}

func RegisterAPICheckRouter(routerType string, f func(fiber.Router)) {
	if _, ok := routerAPICheckMap[routerType]; ok {
		panic("duplicate router type: " + routerType)
	}
	routerAPICheckMap[routerType] = f
}

func RegisterAdminCheckRouter(routerType string, f func(fiber.Router)) {
	if _, ok := routerAdminCheckMap[routerType]; ok {
		panic("duplicate router type: " + routerType)
	}
	routerAdminCheckMap[routerType] = f
}

func RegisterAdminObject(objs []object.AdminObject) {
	adminObjects = append(adminObjects, objs...)
}

func GetAdminObjects() []object.AdminObject {
	sort.SliceStable(adminObjects, func(i, j int) bool {
		return adminObjects[i].Weight < adminObjects[j].Weight
	})
	return adminObjects
}
