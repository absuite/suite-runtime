package main

import (
	"fmt"
	"time"

	"github.com/absuite/suite-runtime/configs"
	"github.com/absuite/suite-runtime/http/middleware"
	"github.com/absuite/suite-runtime/http/routes"
	"github.com/absuite/suite-runtime/repositories"
	"github.com/absuite/suite-runtime/services"
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/utils"
	"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"go.uber.org/dig"
)

func runApp() {
	//创建应用
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.UseGlobal(middleware.Ent)

	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		ctx.JSON(iris.Map{"msg": "Not Found : " + ctx.Path()})
	})
	// 创建容器
	container := dig.New()
	// 基础组件注册
	ComRegister(container, app)
	// 服务注册
	services.Register(container)
	//路由注册
	routes.Register(container)
	//启动服务
	if err := app.Run(iris.Addr(":"+configs.Default.App.Port), iris.WithConfiguration(iris.YAML(utils.JoinCurrentPath("env/iris.yaml")))); err != nil {
		glog.Errorf("Run service error:%s\n", err.Error())
	}
}

func ComRegister(container *dig.Container, app *iris.Application) {
	// 注册app
	if err := container.Provide(func() *iris.Application {
		return app
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
	// 注册db
	config := mysql.Config{User: configs.Default.Db.Username, Passwd: configs.Default.Db.Password, Net: "tcp", Addr: configs.Default.Db.Host, DBName: configs.Default.Db.Database, AllowNativePasswords: true}
	if configs.Default.Db.Port > 0 {
		config.Addr = fmt.Sprintf("%s:%d", configs.Default.Db.Host, configs.Default.Db.Port)
	}
	str := config.FormatDSN()
	orm, err := xorm.NewEngine("mysql", str)
	if err != nil {
		glog.Errorf("orm failed to initialized: %v", err)
	}
	location, err := time.LoadLocation("Asia/Shanghai")
	if err == nil {
		orm.TZLocation = location
	}
	orm.Logger().SetLevel(core.LOG_DEBUG)
	if err := orm.Ping(); err != nil {
		glog.Errorf("can not ping xorm: %v", err)
	}
	iris.RegisterOnInterrupt(func() {
		orm.Close()
	})
	if err := container.Provide(func() *repositories.ModelRepo {
		return repositories.NewModelRepo(orm)
	}); err != nil {
		glog.Errorf("di Provide error:%s", err)
	}
}
