package main

import (
	"os"

	"github.com/absuite/suite-runtime/configs"
	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/utils"
	"github.com/kardianos/service"
)

func main() {
	// 设置日志目录
	glog.SetPath(utils.JoinCurrentPath("storage/logs"))
	// 创建默认配置
	configs.New()
	// 注册服务
	RegisterSv()
}

type program struct {
	cfg *service.Config
}

func (p *program) Start(s service.Service) error {
	glog.Infof("Start  server :%s\n", p.cfg.Name)
	go p.run()
	return nil
}
func (p *program) Stop(s service.Service) error {
	glog.Info("Stop  server:%s\n", p.cfg.Name)
	return nil
}
func (p *program) run() {
	glog.Info("Run  server:%s\n", p.cfg.Name)
	runApp()
}
func RegisterSv() {
	var s = &program{cfg: &service.Config{
		Name:        "SuiteRuntime",
		DisplayName: "Suite Runtime Service",
		Description: "This is a service for suite.",
	}}
	sys := service.ChosenSystem()
	srv, err := sys.New(s, s.cfg)
	if err != nil {
		glog.Fatalf("Init service error:%s\n", err.Error())
	}
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			err := srv.Install()
			if err != nil {
				glog.Fatalf("Install service error:%s\n", err.Error())
			} else {
				glog.Infof("Install service success!")
			}
		case "uninstall":
			err := srv.Uninstall()
			if err != nil {
				glog.Fatalf("Uninstall service error:%s\n", err.Error())
			} else {
				glog.Infof("Uninstall service success!")
			}
		}
		return
	}
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Run programe error:%s\n", err.Error())
	}
}
