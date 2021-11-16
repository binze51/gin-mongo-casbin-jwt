package model

import "gitee.com/binze/binzekeji/pkg/engine"

const (
	TAccount = "bz_accounts"
)

const (
	AesKey = "MQZJlxv3vq0nV7PL"
	AesIv  = "MQZJlxv3vq0nV7PL"
)

type Login struct {
	Account  string `form:"account"  binding:"required"`
	Password string `form:"password"  binding:"required"`
}

//Account 账号实体
type Account struct {
	engine.MODEL
	Name      string `json:"name,omitempty" bson:"name,omitempty"  jsonschema:"required,description=账号名称"`
	Pass      string `json:"password,omitempty" bson:"password,omitempty"  jsonschema:"required,description=账号密码"`
	InitPwd   bool   `json:"init_pwd,omitempty" bson:"init_pwd,omitempty"  jsonschema:"-,description=初始密码标识"`
	Role      string `json:"role,omitempty"  bson:"role,omitempty" jsonschema:"required,description=账号角色"`
	Mobile    string `json:"mobile,omitempty"  bson:"mobile,omitempty" jsonschema:"required,description=账号手机号"`
	Email     string `json:"email,omitempty" bson:"email,omitempty" jsonschema:"-,description=账号邮箱"`
	AvatarUri string `json:"avatar_uri,omitempty" bson:"avatar_uri,omitempty" jsonschema:"-,description=账号头像"`
}

func (a *Account) GetCollName() string {
	return TAccount
}
