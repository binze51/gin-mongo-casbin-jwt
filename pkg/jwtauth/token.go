package jwtauth

import (
	"github.com/dgrijalva/jwt-go"
)

// tokenInfo token信息
type tokenInfo struct {
	AccountName      string `json:"accountName,omitempty"`
	AccountRole      string `json:"accountRole,omitempty"`
	AccountEmail     string `json:"accountEmail,omitempty"`
	AccountMobile    string `json:"accountMobile,omitempty"`
	AccountInitPwd   bool   `json:"initPwd,omitempty"`
	AccountAvatarUri string `json:"accountAvatarUri,omitempty"`
	TokenType        string `json:"tokenType"`   // 令牌类型
	AccessToken      string `json:"accessToken"` // 访问令牌
	ExpiresAt        int64  `json:"expiresAt"`   // 令牌到期时间
	RefreshToken     struct {
		ExpiresAt int64  `json:"expiresAt"`
		Token     string `json:"token"`
	} `json:"refreshToken"`
}

//增加个人标签信息
type claims struct {
	UserID string `json:"userid"`
	Role   string `json:"role"`
	jwt.StandardClaims
}
