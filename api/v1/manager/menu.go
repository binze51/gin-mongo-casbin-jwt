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

func CreateMenu(c *gin.Context) {
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

func GetMenuList(c *gin.Context) {
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
	filter := bson.M{"parentId": 0}
	menucoll := svc.ManagerSvc.Collection(manager.TSys_roles)
	total, err := menucoll.Find(ctx, filter).Count()
	if err != nil {
		svc.ManagerSvc.Log.Error("GetMenusErr", zap.String("MongoERR", fmt.Sprintf("mongo table:(%s) - %s", manager.TSys_roles, err.Error())))
		response.FailWithErrorRaw(model.AccountCreateError, err.Error(), c)
		return
	}
	var menulist []manager.SysMenu
	curr := menucoll.Find(ctx, filter).Limit(int64(limit)).Skip(int64(offset))
	curr.All(&menulist)
	if len(menulist) > 0 {
		for k := range menulist {
			err = findChildrenMenu(&menulist[k])
		}
	}
	if err != nil {
		svc.ManagerSvc.Log.Error("GetRolesErr", zap.String("MongoERR", fmt.Sprintf("mongo table:(%s) - %s", manager.TSys_roles, err.Error())))
		response.FailWithErrorRaw(model.AccountCreateError, err.Error(), c)
		return
	}
	response.OkWithData(model.PageResult{
		List:     menulist,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, c)
}
