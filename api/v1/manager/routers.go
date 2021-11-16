package manager

import (
	svc "gitee.com/binze/binzekeji/service/manager"
	"github.com/gin-gonic/gin"
)

// RegisterAPI 注册模块 api路由
func RegisterRouterRouterGroup(apiRouter *gin.RouterGroup) {
	//开放接口
	pubRouter := apiRouter.Group("/v1")
	{
		pubRouter.GET("healthz", svc.ManagerSvc.Healthz)
		// Use(svc.ManagerSvc.AddApiRow(pubRouter.BasePath()+"/healthz", "健康探针", "get", "健康监测接口"))
		pubRouter.POST("login", svc.ManagerSvc.AddApiRow(pubRouter.BasePath(), "登录", "post", "登录接口"), Login)
	}

	v1 := apiRouter.Group("/v1/manager", svc.ManagerSvc.AuthPromit)

	//用户管理，curd，d软删除，禁用
	user := v1.Group("")
	{
		user.POST("logout", Logout).Use(svc.ManagerSvc.AddApiRow(user.BasePath()+"/logout", "post", "系统管理", "刷新token"))
		user.POST("create", CreateUser)
		user.PUT("update/:name", UpdateUserByName)
		user.PUT("updateEmail/:name", UpdateUserEmailByName)
		user.PUT("updateMobile/:name", UpdateUserMobileByName)
		user.PUT("updateAvatar/:name", UpdateUserAvatarByName)
		user.PUT("changePwd/:name", ChangeUserPwd)
		user.POST("unlock", UnlockUserByName)

		user.POST("uploadImg", svc.ManagerSvc.Upload)
		user.GET("loadImg/:id", svc.ManagerSvc.LoadImgByID)
	}
	// //菜单管理，curd
	// menu := v1.Group("")
	// {
	// 	menu.GET("menus", CreateUser)
	// 	menu.POST("create_menu", CreateUser)
	// }
	// //api管理，ur，api启用状态
	// apis := v1.Group("")
	// {
	// 	apis.GET("apis", CreateUser)
	// }
	// //rbac授权更新,基于role的修改
	// casbin := v1.Group("")
	// {
	// 	casbin.POST("update_casbin")
	// }

	//角色管理,crud
	role := v1.Group("")
	{
		role.POST("create_role")
		role.GET("roles")
		role.DELETE("role/:id")
		role.POST("update_rbac", UpdateRBACPolicys)
	}

}
