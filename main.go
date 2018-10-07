package main

import (
	"fmt"
	"time"

	"github.com/absuite/suite-runtime/http/middleware"
	"github.com/absuite/suite-runtime/http/routes"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"

	"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

func main() {

	if appConfig.Port == "" {
		appConfig.Port = "8080"
	}
	app := iris.New()
	app.Logger().SetLevel("debug")

	app.Use(recover.New())
	app.Use(logger.New())
	app.UseGlobal(middleware.Ent)

	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		ctx.JSON(iris.Map{"msg": "Not Found : " + ctx.Path()})
	})
	orm := handleDb(app)
	//路由注册
	routes.Register(&routes.Context{App: app, Orm: orm})
	//启动服务
	app.Run(iris.Addr(":"+appConfig.Port), iris.WithConfiguration(iris.YAML("./configs/iris.yaml")))
}
func handleDb(app *iris.Application) *xorm.Engine {
	//init db
	//username:password@protocol(address)/dbname?param=value

	config := mysql.Config{User: dbConfig.Username, Passwd: dbConfig.Password, Net: "tcp", Addr: dbConfig.Host, DBName: dbConfig.Database, AllowNativePasswords: true}
	if dbConfig.Port > 0 {
		config.Addr = fmt.Sprintf("%s:%d", dbConfig.Host, dbConfig.Port)
	}
	str := config.FormatDSN()
	orm, err := xorm.NewEngine("mysql", str)
	if err != nil {
		app.Logger().Fatalf("orm failed to initialized: %v", err)
	}
	location, err := time.LoadLocation("Asia/Shanghai")
	orm.TZLocation = location
	// orm.ShowSQL(true)
	orm.Logger().SetLevel(core.LOG_DEBUG)
	if err := orm.Ping(); err != nil {
		app.Logger().Fatalf("can not ping xorm: %v", err)
	}
	iris.RegisterOnInterrupt(func() {
		orm.Close()
	})
	return orm
}
