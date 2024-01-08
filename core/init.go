package core

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var dbMap = map[string]func([]*gorm.DB){}

var routerNoCheckMap = map[string]func(fiber.Router){}
var routerCheckMap = map[string]func(fiber.Router){}
var routerAdminCheckMap = map[string]func(fiber.Router){}

func RegisterDatabase(dbType string, f func([]*gorm.DB)) {
	if _, ok := dbMap[dbType]; ok {
		panic("duplicate db type: " + dbType)
	}
	dbMap[dbType] = f
}

func RegisterNoCheckRouter(routerType string, f func(fiber.Router)) {
	if _, ok := routerNoCheckMap[routerType]; ok {
		panic("duplicate router type: " + routerType)
	}
	routerNoCheckMap[routerType] = f
}

func RegisterCheckRouter(routerType string, f func(fiber.Router)) {
	if _, ok := routerCheckMap[routerType]; ok {
		panic("duplicate router type: " + routerType)
	}
	routerCheckMap[routerType] = f
}

func RegisterAdminCheckRouter(routerType string, f func(fiber.Router)) {
	if _, ok := routerAdminCheckMap[routerType]; ok {
		panic("duplicate router type: " + routerType)
	}
	routerAdminCheckMap[routerType] = f
}
