package manager

import "gitee.com/binze/binzekeji/pkg/engine"

const (
	TSys_roles = "sys_roles"
)

//角色树
type SysRole struct {
	engine.MODEL
	RoleId          uint      `json:"roleid"`                        //角色id，是前端创建的，主键唯一性，也是存到casbin的v0里的值
	Name            string    `json:"authorityName" bson:"name:角色名"` // 角色别名
	ParentId        uint      `json:"parentId" bson:"parentId"`      // 父角色ID
	Children        []SysRole `json:"children" bson:"Children"`      // 数据库 维护上下两级 即可，它的子角色，一对多关系
	DefaultRootMenu string    `json:"defaultRouter" bson:"-"`        // 默认菜单(默认dashboard菜单)
}

//新建新角色时的默认菜单，必须有一个根菜单
func DefaultRootMenu() []SysMenu {
	return []SysMenu{{
		ParentId:  "0",
		Path:      "dashboard",
		Name:      "dashboard",
		Component: "view/dashboard/index.vue",
		Sort:      1,
		Title:     "仪表盘",
		Icon:      "setting",
	}}
}

//获取角色树（所有角色列表）

//角色里有菜单选配，有api选配，这些数据 都要直接获取数据库，不能组装成树
