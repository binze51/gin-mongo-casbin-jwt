package manager

import "gitee.com/binze/binzekeji/pkg/engine"

const (
	TSys_users = "sys_users"
)

type SysUser struct {
	engine.MODEL
	Name      string     `json:"name,omitempty" bson:"name,omitempty"  jsonschema:"required,description=账号名"`
	NickName  string     `json:"nick_name,omitempty" bson:"nick_name,omitempty"  jsonschema:"required,description=账号昵称"`
	Pass      string     `json:"password,omitempty" bson:"password,omitempty"  jsonschema:"required,description=账号密码"`
	InitPwd   bool       `json:"init_pwd,omitempty" bson:"init_pwd,omitempty"  jsonschema:"-,description=初始密码标识"`
	Role      string     `json:"role,omitempty"  bson:"role,omitempty" jsonschema:"required,description=账号角色"`
	Mobile    string     `json:"mobile,omitempty"  bson:"mobile,omitempty" jsonschema:"required,description=账号手机号"`
	Email     string     `json:"email,omitempty" bson:"email,omitempty" jsonschema:"-,description=账号邮箱"`
	AvatarUri string     `json:"avatar_uri,omitempty" bson:"avatar_uri,omitempty" jsonschema:"-,description=账号头像"`
	Roles     []*SysRole `json:"roles,omitempty" bson:"roles,omitempty" jsonschema:"-,description=账号角色集"` //用户携带的角色列表,这里角色 可以只要必要返回就可以了 前端选中使用场景：1:先获取所有角色（id，别名），2: 在所有菜单列表里选中用户数据里存在的即可
}

func (u *SysUser) GetCollName() string {
	return TSys_users
}

const (
	AesKey = "MQZJlxv3vq0nV7PL"
	AesIv  = "MQZJlxv3vq0nV7PL"
)

type Login struct {
	AccountName string `form:"accountName" binding:"required"`
	Password    string `form:"password" binding:"required"`
	Platform    string `form:"platform" binding:"required"` //x-android || x-pc
	OtherLogout bool   `form:"otherLogout"`
}
