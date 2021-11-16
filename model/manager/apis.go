package manager

import "gitee.com/binze/binzekeji/pkg/engine"

const (
	TSys_apis = "sys_apis"
)

type SysApi struct {
	engine.MODEL
	Name   string `json:"name"`   // api中文描述 "用户登录（必选）
	Path   string `json:"path"`   // api路径, "/base/login"  权限判断1
	Group  string `json:"Group"`  // api所属组  "base"
	Method string `json:"method"` // 方法:创建POST(默认)|查看GET|更新PUT|删除DELETE   "POST" //权限判断2

}

//写个中间件在初始化时写入到数据库，这个是全局数据

//全局path里去更新casbin，更新条件：role，path，meth，也是对应casbin的策略
