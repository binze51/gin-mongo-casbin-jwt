# go 后端项目脚手架

简单，高效，最小化必备，模块化设计可以做业务开发脚手架

## 功能（这一版 demo 主要是后台系统接口）

- 账号登录 退出
- 菜单权限，api 权限
- 授权 token，刷新 token
- 登录错误重试限制，强制退出
- 上传文件，更新文件，下载文件 存 mongo gridfs
- 访问限速保护,每秒 200 并发
- h2c httpServer
- rbac 权限控制
- jwt rsa 身份认证
- 创建和更新操作的幂等设计（下订单场景）
  create：资源 id 需要前端提交时携带，这个 id 服务端可以事先生成给前端
  update：加个 version 版本号，订单 id，订单状态，version （期望的 version）

## 服务启动

export envflag=dev && go run main.go

```shell
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /api/v1/healthz           --> gitee.com/binze/binzekeji/pkg/engine.(*Engine).Healthz-fm (5 handlers)
[GIN-debug] POST   /api/v1/login             --> gitee.com/binze/binzekeji/api/v1/manager.Login (6 handlers)
[GIN-debug] POST   /api/v1/manager/logout    --> gitee.com/binze/binzekeji/api/v1/manager.Logout (6 handlers)
[GIN-debug] POST   /api/v1/manager/create    --> gitee.com/binze/binzekeji/api/v1/manager.CreateUser (7 handlers)
[GIN-debug] PUT    /api/v1/manager/update/:name --> gitee.com/binze/binzekeji/api/v1/manager.UpdateUserByName (7 handlers)
[GIN-debug] PUT    /api/v1/manager/updateEmail/:name --> gitee.com/binze/binzekeji/api/v1/manager.UpdateUserEmailByName (7 handlers)
[GIN-debug] PUT    /api/v1/manager/updateMobile/:name --> gitee.com/binze/binzekeji/api/v1/manager.UpdateUserMobileByName (7 handlers)
[GIN-debug] PUT    /api/v1/manager/updateAvatar/:name --> gitee.com/binze/binzekeji/api/v1/manager.UpdateUserAvatarByName (7 handlers)
[GIN-debug] PUT    /api/v1/manager/changePwd/:name --> gitee.com/binze/binzekeji/api/v1/manager.ChangeUserPwd (7 handlers)
[GIN-debug] POST   /api/v1/manager/unlock    --> gitee.com/binze/binzekeji/api/v1/manager.UnlockUserByName (7 handlers)
[GIN-debug] POST   /api/v1/manager/uploadImg --> gitee.com/binze/binzekeji/service/manager.(*ManagerService).Upload-fm (7 handlers)
[GIN-debug] GET    /api/v1/manager/loadImg/:id --> gitee.com/binze/binzekeji/service/manager.(*ManagerService).LoadImgByID-fm (7 handlers)
[GIN-debug] POST   /api/v1/manager/create_role --> gitee.com/binze/binzekeji/pkg/engine.(*Engine).AuthPromit-fm (5 handlers)
[GIN-debug] GET    /api/v1/manager/roles     --> gitee.com/binze/binzekeji/pkg/engine.(*Engine).AuthPromit-fm (5 handlers)
[GIN-debug] DELETE /api/v1/manager/role/:id  --> gitee.com/binze/binzekeji/pkg/engine.(*Engine).AuthPromit-fm (5 handlers)
[GIN-debug] POST   /api/v1/manager/update_rbac --> gitee.com/binze/binzekeji/api/v1/manager.UpdateRBACPolicys (6 handlers)
[GIN-debug] Listening and serving H2C on 6010
[GIN-debug] RATE   Current Rate Limit: 200 requests/s
```
