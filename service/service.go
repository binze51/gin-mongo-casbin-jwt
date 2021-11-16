package service

import (
	"gitee.com/binze/binzekeji/config"
	"gitee.com/binze/binzekeji/pkg/engine"
	"gitee.com/binze/binzekeji/service/manager"
)

//要仅执行一次
func InitService(e *engine.Engine, conf *config.Config) {
	manager.NewManagerSvc(e, conf)
	// CustomerSvc = manager.NewManagerSvc(e)

}
