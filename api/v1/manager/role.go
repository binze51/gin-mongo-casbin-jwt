package manager

import (
	"context"
	"fmt"

	"gitee.com/binze/binzekeji/model"
	"gitee.com/binze/binzekeji/model/manager"
	"gitee.com/binze/binzekeji/pkg/response"
	"gitee.com/binze/binzekeji/pkg/utils"
	svc "gitee.com/binze/binzekeji/service/manager"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func CreateRole(c *gin.Context) {
	postInfo := new(createRequest)

	if err := c.ShouldBind(postInfo); err != nil {
		svc.ManagerSvc.Log.Error("CreateAccountErr", zap.String("MongoERR", err.Error()))
		response.Fail(response.ParamBindError, c)
		return
	}
	aes := utils.NewAes(model.AesKey, model.AesIv)
	encryptPwd, _ := aes.Encrypt(postInfo.AccountPwd)
	postInfo.AccountPwd = encryptPwd
	_, err := svc.ManagerSvc.Collection(model.TAccount).InsertOne(context.TODO(), postInfo)
	if err != nil {
		svc.ManagerSvc.Log.Error("CreateAccountErr", zap.String("MongoERR", fmt.Sprintf("mongo table:(%s) - %s", model.TAccount, err.Error())))
		response.FailWithErrorRaw(model.AccountCreateError, err.Error(), c)
		return
	}
	response.Ok(c)

}

func DeleteRole(c *gin.Context) {
	//删除p的所有策略数据
	//全量重新新增
}
func GetRoleList(c *gin.Context) {
	var pageInfo model.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		svc.ManagerSvc.Log.Error("GetMenusErr", zap.String("MongoERR", fmt.Sprintf("mongo table:(%s) - %s", manager.TSys_roles, err.Error())))
		response.FailWithErrorRaw(model.AccountCreateError, err.Error(), c)
		return
	}
	ctx := context.TODO()
	limit := pageInfo.PageSize
	offset := pageInfo.PageSize * (pageInfo.Page - 1)
	filter := bson.M{"parentId": "00"}
	total, err := svc.ManagerSvc.Collection(manager.TSys_roles).Find(ctx, filter).Count()
	if err != nil {
		svc.ManagerSvc.Log.Error("GetRolesErr", zap.String("MongoERR", fmt.Sprintf("mongo table:(%s) - %s", manager.TSys_roles, err.Error())))
		response.FailWithErrorRaw(model.AccountCreateError, err.Error(), c)
		return
	}
	var rolelist []manager.SysRole
	curr := svc.ManagerSvc.Collection(manager.TSys_roles).Find(ctx, filter).Limit(int64(limit)).Skip(int64(offset))
	curr.All(&rolelist)
	if len(rolelist) > 0 {
		for k := range rolelist {
			err = findChildrenRole(&rolelist[k])
		}
	}
	if err != nil {
		svc.ManagerSvc.Log.Error("GetRolesErr", zap.String("MongoERR", fmt.Sprintf("mongo table:(%s) - %s", manager.TSys_roles, err.Error())))
		response.FailWithErrorRaw(model.AccountCreateError, err.Error(), c)
		return
	}
	response.OkWithData(model.PageResult{
		List:     rolelist,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, c)
}

func UpdateRBACPolicys(c *gin.Context) {
	info := new(manager.CasbinReq)
	_ = c.ShouldBindJSON(info)
	Policies := [][]string{}
	for _, v := range info.CasbinInfos {
		Policies = append(Policies, []string{info.RoleID, v.Path, v.Method})
	}
	success, err := svc.ManagerSvc.Ce.AddPolicies(Policies)
	if !success {
		response.FailWithErrorRaw(model.AccountCreateError, err.Error(), c)
		return
	}
	response.Ok(c)
}

//添加rbac策略
// func (e *Engine) UpdateRBACPolicyByID(p, path, meth string) gin.HandlerFunc {
// 	if ok, _ := e.Ce.AddPolicy(p, path, meth); !ok {
// 		//存在则忽略
// 		log.Println("the policy already exists")
// 	} else {
// 		//添加ok
// 		log.Println("add one policy successfully")
// 	}
// 	return func(c *gin.Context) {
// 		c.Next()
// 	}
// }

// }
// casbinRouter.POST("/updateall", _engine.AddRBACPolicy)
