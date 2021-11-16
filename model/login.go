package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	TLoginList = "bz_loginlist"
)

//LoginList 登录列表
type LoginList struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Platform     string             `bson:"platform" `
	AccountName  string             `bson:"accountName"`
	Lock         bool               `bson:"lock,omitempty"`
	CreateAt     int64              `bson:"createAt"`
	AccessToken  string             `bson:"accessToken"`
	RefreshToken string             `bson:"refreshToken,omitempty"`
	Errs         uint               `bson:"errs,omitempty"` //错误次数，最大约束 5次
	Ttl          time.Time          `bson:"ttl,omitempty"`  //refreshToken ttl 24h
}
