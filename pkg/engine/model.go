package engine

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//G_MODEL 实体通用结构
type MODEL struct {
	ID       primitive.ObjectID `json:"id,omitempty"  bson:"_id,omitempty" jsonschema:"-"` // 实体唯一标识
	Version  uint               `json:"version,omitempty" bson:"version,omitempty"`        // 实体版本号,实体修改幂等
	CreateAt int64              `json:"create_at,omitempty" bson:"create_at,omitempty"`    // 创建时间
	DeleteAt int64              `json:"delete_at,omitempty" bson:"delete_at,omitempty"`    // 删除时间
	UpdateAt int64              `json:"update_at,omitempty" bson:"update_at,omitempty"`    // 更新时间

	IsDelete  *bool `json:"is_delete,omitempty" bson:"is_delete,omitempty"`   // 是否删除
	IsDisable *bool `json:"is_disable,omitempty" bson:"is_disable,omitempty"` // 是否禁用
}

// Response  通用返回结构
type Response struct {
	BizCode int         `json:"bizcode"` //业务码
	Message string      `json:"message"` //消息特征 success，fail
	Data    interface{} `json:"data"`    //数据
}
