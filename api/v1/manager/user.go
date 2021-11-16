package manager

import (
	"context"
	"fmt"
	"time"

	"gitee.com/binze/binzekeji/model"
	"gitee.com/binze/binzekeji/model/manager"
	"gitee.com/binze/binzekeji/pkg/response"
	"gitee.com/binze/binzekeji/pkg/utils"
	svc "gitee.com/binze/binzekeji/service/manager"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func CreateUser(c *gin.Context) {
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

type createRequest struct {
	AccountName string `form:"account" binding:"required" bson:"accountName,omitempty"`
	AccountPwd  string `form:"password" binding:"required" bson:"password,omitempty"`
	AccountRole string `form:"role,omitempty" bson:"role,omitempty" `
	Mobile      string `form:"mobile,omitempty" bson:"mobile,omitempty"`
	Email       string `form:"email,omitempty" bson:"email,omitempty"`
	AvatarUri   string `form:"avatarUri,omitempty" bson:"avatarUri,omitempty"`
}

type updateRequest struct {
	Mobile    string `form:"mobile,omitempty" `
	Email     string `form:"email,omitempty" `
	AvatarUri string `form:"avatarUri,omitempty"`
}

func UpdateUserByName(c *gin.Context) {
	accountName := c.Param("name")
	if accountName == "" {
		svc.ManagerSvc.Log.Error("UpdateByNameErr", zap.String("ParamERR", "name路径参数缺失"))
		response.Fail(response.ParamBindError, c)
		return
	}
	curAccountName := c.GetString("userid")
	if curAccountName != "admin" && accountName == "admin" {
		svc.ManagerSvc.Log.Error("UpdateEmailByName", zap.String("ERR", "非法修改admin资料"))
		response.Fail(model.NoAdminUpdateAdminError, c)
		return
	}
	account := new(model.Account)
	err := svc.ManagerSvc.Collection(model.TAccount).Find(context.Background(),
		bson.M{"accountName": accountName}).One(account)
	if err != nil {
		svc.ManagerSvc.Log.Error("UpdateByNameErr", zap.String("MongoERR", err.Error()))
		response.Fail(model.AccountSearchError, c)
		return
	}
	info := new(updateRequest)
	err = c.ShouldBind(info)
	if err != nil {
		svc.ManagerSvc.Log.Error("UpdateByNameErr", zap.String("MongoERR", err.Error()))
		response.Fail(response.ParamBindError, c)
		return
	}

	setObj := bson.M{}
	if info.AvatarUri != "" {
		setObj["avatarUri"] = info.AvatarUri
	}
	if info.Email != "" {
		setObj["email"] = info.Email
	}
	if info.Mobile != "" {
		setObj["mobile"] = info.Mobile
	}
	err = svc.ManagerSvc.Collection(model.TAccount).UpdateOne(context.Background(),
		bson.M{"accountName": accountName},
		bson.M{"$set": setObj})
	if err != nil {
		svc.ManagerSvc.Log.Error("UpdateByNameErr", zap.String("MongoERR", err.Error()))
		response.Fail(model.AccountUpdateSaveError, c)
		return
	}
	response.Ok(c)
}

func UpdateUserEmailByName(c *gin.Context) {
	accountName := c.Param("name")
	if accountName == "" {
		svc.ManagerSvc.Log.Error("UpdateEmailByName", zap.String("ParamERR", "name路径参数缺失"))
		response.Fail(response.ParamBindError, c)
		return
	}
	curAccountName := c.GetString("userid")
	if curAccountName != "admin" && accountName == "admin" {
		svc.ManagerSvc.Log.Error("UpdateEmailByName", zap.String("ERR", "非法修改admin资料"))
		response.Fail(model.NoAdminUpdateAdminError, c)
		return
	}
	account := new(model.Account)
	err := svc.ManagerSvc.Collection(model.TAccount).Find(context.Background(),
		bson.M{"accountName": accountName}).One(account)
	if err != nil {
		svc.ManagerSvc.Log.Error("UpdateEmailByName", zap.String("MongoERR", err.Error()))
		response.Fail(model.AccountSearchError, c)
		return
	}
	info := new(updateRequest)
	err = c.ShouldBind(info)
	if err != nil {
		svc.ManagerSvc.Log.Error("UpdateEmailByNameErr", zap.String("ParamERR", err.Error()))
		response.Fail(response.ParamBindError, c)
		return
	}
	if info.Email == "" {
		svc.ManagerSvc.Log.Error("UpdateEmailByNameErr", zap.String("ParmERR", "邮箱参数缺失"))
		response.Fail(model.UpdateEmailError, c)
		return
	}
	err = svc.ManagerSvc.Collection(model.TAccount).UpdateOne(context.Background(),
		bson.M{"accountName": accountName},
		bson.M{"$set": bson.M{"email": info.Email}})
	if err != nil {
		svc.ManagerSvc.Log.Error("UpdateEmailByNameErr", zap.String("MongoERR", err.Error()))
		response.Fail(model.UpdateEmailError, c)
		return
	}
	response.Ok(c)
}
func UpdateUserMobileByName(c *gin.Context) {
	accountName := c.Param("name")
	if accountName == "" {
		svc.ManagerSvc.Log.Error("UpdateMobileByNameErr", zap.String("ParamERR", "name路径参数缺失"))
		response.Fail(response.ParamBindError, c)
		return
	}
	curAccountName := c.GetString("userid")
	if curAccountName != "admin" && accountName == "admin" {
		svc.ManagerSvc.Log.Error("UpdateEmailByName", zap.String("ERR", "非法修改admin资料"))
		response.Fail(model.NoAdminUpdateAdminError, c)
		return
	}
	account := new(model.Account)
	err := svc.ManagerSvc.Collection(model.TAccount).Find(context.Background(),
		bson.M{"accountName": accountName}).One(account)
	if err != nil {
		svc.ManagerSvc.Log.Error("UpdateMobileByNameErr", zap.String("MongoERR", err.Error()))
		response.Fail(model.AccountSearchError, c)
		return
	}
	info := new(updateRequest)
	err = c.ShouldBind(info)
	if err != nil {
		svc.ManagerSvc.Log.Error("UpdateMobileByNameErr", zap.String("ParamERR", err.Error()))
		response.Fail(response.ParamBindError, c)
		return
	}
	if info.Mobile == "" {
		svc.ManagerSvc.Log.Error("UpdateEmailByNameErr", zap.String("ParmERR", "手机号参数缺失"))
		response.Fail(model.UpdateMobileError, c)
		return
	}
	err = svc.ManagerSvc.Collection(model.TAccount).UpdateOne(context.Background(),
		bson.M{"accountName": accountName},
		bson.M{"$set": bson.M{"mobile": info.Mobile}})
	if err != nil {
		svc.ManagerSvc.Log.Error("UpdateByNameErr", zap.String("MongoERR", err.Error()))
		response.Fail(model.UpdateMobileError, c)
		return
	}
	response.Ok(c)
}
func UpdateUserAvatarByName(c *gin.Context) {
	accountName := c.Param("name")
	if accountName == "" {
		svc.ManagerSvc.Log.Error("UpdateAvatarByNameErr", zap.String("ParamERR", "name路径参数缺失"))
		response.Fail(response.ParamBindError, c)
		return
	}
	curAccountName := c.GetString("userid")
	if curAccountName != "admin" && accountName == "admin" {
		svc.ManagerSvc.Log.Error("UpdateEmailByName", zap.String("ERR", "非法修改admin资料"))
		response.Fail(model.NoAdminUpdateAdminError, c)
		return
	}
	account := new(model.Account)
	err := svc.ManagerSvc.Collection(model.TAccount).Find(context.Background(),
		bson.M{"accountName": accountName}).One(account)
	if err != nil {
		svc.ManagerSvc.Log.Error("UpdateAvatarByNameErr", zap.String("MongoERR", err.Error()))
		response.Fail(model.AccountSearchError, c)
		return
	}
	info := new(updateRequest)
	err = c.ShouldBind(info)
	if err != nil {
		svc.ManagerSvc.Log.Error("UpdateAvatarByNameErr", zap.String("ParamERR", err.Error()))
		response.Fail(response.ParamBindError, c)
		return
	}

	if info.AvatarUri == "" {
		svc.ManagerSvc.Log.Error("UpdateAvatarByNameErr", zap.String("ParmERR", "头像uri参数缺失"))
		response.Fail(model.UpdateAvatarError, c)
		return
	}
	err = svc.ManagerSvc.Collection(model.TAccount).UpdateOne(context.Background(),
		bson.M{"accountName": accountName},
		bson.M{"$set": bson.M{"avatarUri": info.AvatarUri}})
	if err != nil {
		svc.ManagerSvc.Log.Error("UpdateAvatarByNameErr", zap.String("MongoERR", err.Error()))
		response.Fail(model.UpdateAvatarError, c)
		return
	}
	response.Ok(c)
}

type UnlockRequest struct {
	AccountName string `form:"accountName" binding:"required"`
	Platform    string `form:"platform" binding:"required"`
}

func UnlockUserByName(c *gin.Context) {
	UnlockPost := new(UnlockRequest)
	err := c.ShouldBind(UnlockPost)
	if err != nil {
		svc.ManagerSvc.Log.Error("UnlockByNameErr", zap.String("ParamERR", err.Error()))
		response.Fail(response.ParamBindError, c)
		return
	}
	err = svc.ManagerSvc.Collection(model.TDenyList).Remove(context.Background(),
		bson.M{"accountName": UnlockPost.AccountName, "platform": UnlockPost.Platform})
	if err != nil {
		svc.ManagerSvc.Log.Error("UnlockByNameErr", zap.String("MongoERR", err.Error()))
		response.Fail(model.AdminUnlockError, c)
		return
	}
	response.Ok(c)
}

type changePwdRequest struct {
	OldPwd string `form:"oldPwd,omitempty" binding:"required"`
	NewPwd string `form:"newPwd,omitempty" binding:"required"`
}

func ChangeUserPwd(c *gin.Context) {
	accountName := c.Param("name")
	if accountName == "" {
		svc.ManagerSvc.Log.Error("UpdateEmailByName", zap.String("ParamERR", "name路径参数缺失"))
		response.Fail(response.ParamBindError, c)
		return
	}
	curAccountName := c.GetString("userid")
	if curAccountName != "admin" && accountName == "admin" {
		svc.ManagerSvc.Log.Error("UpdateEmailByName", zap.String("ERR", "非法修改admin资料"))
		response.Fail(model.NoAdminUpdateAdminError, c)
		return
	}
	postInfo := new(changePwdRequest)
	err := c.ShouldBind(postInfo)
	if err != nil {
		svc.ManagerSvc.Log.Error("ChangePwdErr", zap.String("ParamERR", err.Error()))
		response.Fail(response.ParamBindError, c)
		return
	}

	account := new(model.Account)
	err = svc.ManagerSvc.Collection(model.TAccount).Find(context.Background(),
		bson.M{"accountName": accountName}).One(account)
	if err != nil {
		svc.ManagerSvc.Log.Error("ChangePwdErr", zap.String("MongoERR", err.Error()))
		response.Fail(model.AccountSearchError, c)
		return
	}
	aes := utils.NewAes(model.AesKey, model.AesIv)
	encryptOldPwd, _ := aes.Encrypt(postInfo.OldPwd)
	if encryptOldPwd != account.Pass {
		svc.ManagerSvc.Log.Error("ChangePwdErr", zap.String("PwdERR", "密码错误"))
		response.Fail(model.AccountChangePwdError, c)
		return
	}
	encryptNewPwd, _ := aes.Encrypt(postInfo.NewPwd)
	setObj := bson.M{"password": encryptNewPwd}
	if accountName == "admin" && account.InitPwd {
		setObj["initPwd"] = false
	}
	err = svc.ManagerSvc.Collection(model.TAccount).UpdateOne(context.Background(),
		bson.M{"accountName": accountName},
		bson.M{"$set": setObj})
	if err != nil {
		svc.ManagerSvc.Log.Error("ChangePwdErr", zap.String("MongoERR", err.Error()))
		response.Fail(model.AccountNewPwdSetFailError, c)
		return
	}
	response.Ok(c)
}

// Account login structure

//Login 登录比对密码，正常签发token，异常：目前直接作废冲突的token，后续可以走服务push给指定客户端，客户端主动发退出请求
func Login(c *gin.Context) {
	loginPost := new(manager.Login)
	if err := c.ShouldBind(loginPost); err != nil {
		svc.ManagerSvc.Log.Error("LoginErr", zap.String("ParamERR", err.Error()))
		response.Fail(response.ParamBindError, c)
		return
	}
	account := new(model.Account)
	AccountColl := svc.ManagerSvc.Collection(model.TAccount)
	err := AccountColl.Find(context.Background(),
		bson.M{"accountName": loginPost.AccountName}).One(account)
	if err != nil {
		curCode, message := checkLoginErrorlimit(loginPost, svc.ManagerSvc.GetDatabase())
		svc.ManagerSvc.Log.Error("LoginErr", zap.String("MongoERR", message))
		response.FailWithError(curCode, message, c)
		return
	}

	aes := utils.NewAes(model.AesKey, model.AesIv)
	encryptpwd, _ := aes.Encrypt(loginPost.Password)
	if encryptpwd != account.Pass {
		curCode, message := checkLoginErrorlimit(loginPost, svc.ManagerSvc.GetDatabase())
		svc.ManagerSvc.Log.Error("LoginErr", zap.String("MongoERR", message))
		response.FailWithError(curCode, message, c)
		return
	}
	newToken, err := svc.ManagerSvc.Jwt.IssueToken(context.Background(),
		false, loginPost.AccountName, account.Role, account.Mobile, account.Email,
		account.AvatarUri, "scheduling-manager", account.InitPwd)
	if err != nil {
		svc.ManagerSvc.Log.Error("LoginErr", zap.String("MongoERR", err.Error()))
		response.Fail(model.IssueTokenFaildError, c)
		return
	}
	if loginPost.OtherLogout {
		//服务端处理：强制踢掉其他人-把LoginList里的记录mv到denytoken里即可
		otherLogin := new(model.LoginList)
		err = svc.ManagerSvc.Collection(model.TLoginList).Find(context.Background(), bson.M{"accountName": loginPost.AccountName,
			"lock": true, "platform": loginPost.Platform}).One(otherLogin)
		if err != nil {
			svc.ManagerSvc.Log.Error("OtherLogoutErr", zap.String("MongoERR", err.Error()))
			response.Fail(model.AccountSearchError, c)
			return
		}
		denyToken := &model.DenyTokenList{
			AccountName: otherLogin.AccountName,
			UseToken:    otherLogin.AccessToken,
			CreateAt:    time.Now().Local().Unix(),
		}
		_, err = svc.ManagerSvc.Collection(model.TDenyTokenList).InsertOne(context.Background(),
			denyToken)
		if err != nil {
			svc.ManagerSvc.Log.Error("OtherLogoutErr", zap.String("MongoERR", err.Error()))
			response.Fail(model.OtherRepeatTokenMoveError, c)
			return
		}
		err = svc.ManagerSvc.Collection(model.TLoginList).Remove(context.Background(),
			bson.M{"_id": otherLogin.ID})
		if err != nil {
			svc.ManagerSvc.Log.Error("OtherLogoutErr", zap.String("MongoERR", err.Error()))
			response.Fail(model.OtherRepeatRefreshTokeClearError, c)
			return
		}
		//客户端处理：发送事件给客户端，他前端程序自动退出
	}
	loginSave := &model.LoginList{
		Platform:     loginPost.Platform,
		AccountName:  loginPost.AccountName,
		Lock:         true,
		RefreshToken: newToken.RefreshToken.Token,
		AccessToken:  newToken.AccessToken,
		CreateAt:     newToken.ExpiresAt,
		Ttl:          time.Now().Local(),
	}
	_, _ = svc.ManagerSvc.Collection(model.TLoginList).RemoveAll(context.Background(),
		bson.M{"accountName": loginPost.AccountName, "errs": bson.M{"$gte": 1}})
	_, err = svc.ManagerSvc.Collection(model.TLoginList).InsertOne(context.Background(),
		loginSave)
	if err != nil {
		svc.ManagerSvc.Log.Error("LoginErr", zap.String("MongoERR", err.Error()))
		response.Fail(model.AccountPlatformLoginRepeatError, c)
		return
	}
	response.OkWithData(newToken, c)
}

// checkLoginErrorlimit 检测错误登录情况，按需求最大5次错误
func checkLoginErrorlimit(inlogin *manager.Login, c *qmgo.Database) (response.RetCode, string) {
	historyLogin := new(model.LoginList)
	loginListColl := c.Collection(model.TAccount)
	err := loginListColl.Find(context.Background(),
		bson.M{"accountName": inlogin.AccountName, "errs": 1}).One(historyLogin)
	if err == nil {
		return model.AccountLoginDenyError,
			model.AccountCodeTable[model.AccountLoginDenyError]
	}
	err = loginListColl.Find(context.Background(),
		bson.M{"accountName": inlogin.AccountName, "errs": bson.M{"$gte": 2}}).One(historyLogin)
	if err != nil {
		errLogin := new(model.LoginList)
		errLogin.AccountName = inlogin.AccountName
		errLogin.CreateAt = time.Now().Local().Unix()
		errLogin.Platform = inlogin.Platform
		errLogin.Errs = 5
		_, err := loginListColl.InsertOne(context.Background(), errLogin)
		if err != nil {
			return model.AccountErrorLoginSaveError, err.Error()
		}
		return model.AccountLoginLimitError,
			fmt.Sprintf(model.AccountCodeTable[model.AccountLoginLimitError], errLogin.Errs-1)
	} else if historyLogin.Errs == 2 {
		denyColl := c.Collection(model.TDenyList)
		_, err = denyColl.InsertOne(context.Background(), bson.M{"platform": inlogin.Platform,
			"accountName": inlogin.AccountName, "createAt": time.Now().Local().Unix(), "ttl": time.Now().Local()})
		if err != nil {
			return model.LoginMoveDenyError,
				model.AccountCodeTable[model.LoginMoveDenyError]
		}
		err = loginListColl.UpdateOne(context.Background(),
			bson.M{"accountName": inlogin.AccountName, "errs": bson.M{"$gte": 1}},
			bson.M{"$set": bson.M{"errs": 1}})
		if err != nil {
			return model.AccountErrorLoginCountError, err.Error()
		}
		return model.AccountLoginDenyError,
			model.AccountCodeTable[model.AccountLoginDenyError]
	} else if historyLogin.Errs > 1 {
		err := loginListColl.UpdateOne(context.Background(),
			bson.M{"accountName": inlogin.AccountName, "errs": bson.M{"$gte": 2}}, bson.M{"$inc": bson.M{"errs": -1}})
		if err != nil {
			return model.AccountErrorLoginCountError, err.Error()
		}
		return model.AccountLoginLimitError,
			fmt.Sprintf(model.AccountCodeTable[model.AccountLoginLimitError], historyLogin.Errs-2)
	} else {
		return model.AccountLoginDenyError,
			model.AccountCodeTable[model.AccountLoginDenyError]
	}
}

// Logout 退出
func Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")[len("Bearer "):]
	uid := c.GetString("userid")

	denyToken := &model.DenyTokenList{
		AccountName: uid,
		UseToken:    token,
		CreateAt:    time.Now().Local().Unix(),
	}
	denyTokenColl := svc.ManagerSvc.Collection(model.TDenyTokenList)
	_, err := denyTokenColl.InsertOne(context.Background(), denyToken)
	if err != nil {
		svc.ManagerSvc.Log.Error("LogoutErr", zap.String("MongoERR", err.Error()))
		response.Fail(model.AccountLogoutError, c)
		return
	}
	response.Ok(c)
}

func RefreshToken(c *gin.Context) {
	accountName := c.GetString("userid")
	role := c.GetString("role")
	refreshToken := new(model.RefreshToken)
	err := c.ShouldBind(refreshToken)
	if err != nil {
		svc.ManagerSvc.Log.Error("RefreshTokenErr", zap.String("MongoERR", err.Error()))
		response.Fail(response.ParamBindError, c)
		return
	}
	loginRet := new(model.LoginList)
	_ = svc.ManagerSvc.Collection(model.TLoginList).Find(context.Background(),
		bson.M{"accountNameErr": accountName, "refreshToken": refreshToken.Refresh}).One(loginRet)
	if loginRet == nil {
		svc.ManagerSvc.Log.Error("RefreshTokenErr", zap.String("ERR", "刷新token已过期"))
		response.Fail(model.RefreshTokenExpiredError, c)
		return
	}
	newToken, err := svc.ManagerSvc.Jwt.IssueToken(context.Background(),
		true, accountName, role, "", "", "", "scheduling-manager", false)
	if err != nil {
		svc.ManagerSvc.Log.Error("RefreshTokenErr", zap.String("MongoERR", err.Error()))
		response.Fail(model.IssueTokenFaildError, c)
		return
	}
	response.OkWithData(newToken, c)
}
