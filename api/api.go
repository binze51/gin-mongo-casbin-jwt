package api

import (
	"gitee.com/binze/binzekeji/api/v1/manager"
	"gitee.com/binze/binzekeji/pkg/router"
)

func RegisterRouter(rg *router.Router) {
	apigroute := rg.Group("/api")
	//后台接口
	manager.RegisterRouterRouterGroup(apigroute)

}
