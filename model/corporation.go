package model

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	TCorporation = "bz_corporations"
)

//Corporation 客户企业信息
type Corporation struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name    string             `form:"name" json:"name,omitempty" bson:"name,omitempty"`
	Address string             `form:"address" json:"address,omitempty" bson:"address,omitempty"`
	Desc    string             `form:"desc" json:"desc,omitempty" bson:"desc,omitempty"`
}
