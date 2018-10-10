package main

import (
	"fmt"
	"time"

	"github.com/absuite/suite-runtime/configs"
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

	configs.New()
	app := iris.New()

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
	app.Run(iris.Addr(":"+configs.Default.App.Port), iris.WithConfiguration(iris.YAML("./env/iris.yaml")))
}
func handleDb(app *iris.Application) *xorm.Engine {
	//init db
	//username:password@protocol(address)/dbname?param=value

	config := mysql.Config{User: configs.Default.Db.Username, Passwd: configs.Default.Db.Password, Net: "tcp", Addr: configs.Default.Db.Host, DBName: configs.Default.Db.Database, AllowNativePasswords: true}
	if configs.Default.Db.Port > 0 {
		config.Addr = fmt.Sprintf("%s:%d", configs.Default.Db.Host, configs.Default.Db.Port)
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
