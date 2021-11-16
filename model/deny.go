package model

import "time"

const (
	TDenyList      = "bz_denylist"
	TDenyTokenList = "bz_denytokenlist"
)

//Denylist uid+Platform黑名单
type DenyList struct {
	Platform    string    `bson:"platform,omitempty"`
	AccountName string    `bson:"accountName,omitempty"`
	CreateAt    int64     `bson:"createAt,omitempty"`
	Ttl         time.Time `bson:"ttl,omitempty"` //12h ttl或管理清除解锁
}

//DenyTokenList token黑名单
type DenyTokenList struct {
	AccountName string `bson:"AccountName,omitempty"`
	UseToken    string `bson:"useToken,omitempty"`
	CreateAt    int64  `bson:"createAt,omitempty"`
}
