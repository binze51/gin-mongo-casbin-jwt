package manager

import "gitee.com/binze/binzekeji/pkg/engine"

const (
	TSys_menus = "sys_menus"
)

//菜单树
type SysMenu struct {
	engine.MODEL
	ID       uint      `json:"菜单"`                            //菜单id
	Path     string    `json:"path" gorm:"comment:路由path"`    // 路由path
	Name     string    `json:"name" gorm:"comment:路由name"`    // 路由name
	ParentId string    `json:"parentId" gorm:"comment:父菜单ID"` // 父菜单ID ParentId:0是一项一级菜单
	Children []SysMenu `json:"children" bson:"children"`      //一对多关系

	Hidden    bool   `json:"hidden" gorm:"comment:是否在列表隐藏"`     // 是否在列表隐藏
	Component string `json:"component" gorm:"comment:对应前端文件路径"` // 对应前端文件路径
	Sort      int    `json:"sort" gorm:"comment:排序标记"`          // 排序标记
	Title     string `json:"title" gorm:"comment:菜单名"`          // 菜单名
	Icon      string `json:"icon" gorm:"comment:菜单图标"`          // 菜单图标
	MenuLevel uint   `json:"-"`
}

// 为角色增加menu树--写到数据库表里role_menu
// 为角色增加api集合--更新角色的casbin 策略记录

//获取菜单树（获取所有菜单列表）
