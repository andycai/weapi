package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
)

type Querier interface {
	// SELECT * FROM @@table WHERE name = @name{{if role !=""}} AND role = @role{{end}}
	FilterWithNameAndRole(name, role string) ([]gen.T, error)
}

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./v2/dao",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	// gormdb, _ := gorm.Open(mysql.Open("root:123456@(127.0.0.1:3306)/werite?charset=utf8mb4&parseTime=true&loc=Local"))
	gormdb, _ := gorm.Open(sqlite.Open("./werite.db"))
	g.UseDB(gormdb)

	g.ApplyBasic(g.GenerateAllTable()...)

	//g.ApplyInterface(func(querier Querier) {}, model.User{}, model.Group{}, model.Activity{})

	// post := g.GenerateModel("post", gen.FieldRelateModel(field.HasMany, "Tags", model.Tag{},
	// 	&field.RelateConfig{
	// 		// RelateSlice: true,
	// 		GORMTag: field.GormTag{"foreignKey": []string{"CustomerRefer"}, "references": []string{"ID"}},
	// 	}),
	// )
	// g.ApplyBasic(post)

	g.Execute()
}
