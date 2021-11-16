package jwtauth

import (
	"context"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 定义错误
var (
	ErrInvalidToken = errors.New("invalid token")
)

type JwtConfig struct {
	TokenType        string `yaml:"tokentype"`
	Expired          int64  `yaml:"expired"`        //24h
	RefreshExpired   int64  `yaml:"refreshexpired"` //7天
	RSAPrivateSecret string `yaml:"privatekey"`
	RSAPublicSecret  string `yaml:"publickey"`
}

// New 创建认证实例
func New(conf *JwtConfig) *JWTAuth {
	return &JWTAuth{
		conf: conf,
	}
}

// JWTAuth jwt认证
type JWTAuth struct {
	conf *JwtConfig
}

// GenerateToken 签发令牌
func (a *JWTAuth) IssueToken(ctx context.Context, refresh bool, userID, role, mobile, email, avatarUri, subject string, initPwd bool) (*tokenInfo, error) {
	now := time.Now().Local()
	expiresAt := now.Add(time.Duration(a.conf.Expired) * time.Second).Unix()
	claims := claims{
		userID,
		role,
		jwt.StandardClaims{
			Issuer:    "gosvc.com", //颁发者
			Subject:   subject,     //该token使用实体
			IssuedAt:  now.Unix(),  //颁发时间
			ExpiresAt: expiresAt,   //token过期时间，目前是24h
			NotBefore: now.Unix(),  //token生效时间
		},
	}
	//使用非对称RS256算法来生成JWT
	//生成前两段token结构
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	//这个值需要json化
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(a.conf.RSAPrivateSecret))
	if err != nil {
		return nil, err
	}
	privJwk, err := jwk.New(privateKey)
	if err != nil {
		return nil, err
	}
	// generates Kid using Key.Thumbprint method with crypto.SHA256
	jwk.AssignKeyID(privJwk) //nolint:errcheck
	token.Header["kid"] = privJwk.KeyID()
	//返回JWT的token字符串
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return nil, err
	}
	//token和刷新token的信息
	tokenInfo := new(tokenInfo)
	if refresh {
		tokenInfo.TokenType = "Bearer"
		tokenInfo.ExpiresAt = expiresAt
		tokenInfo.AccessToken = tokenString
	} else {
		tokenInfo.AccountName = userID
		tokenInfo.AccountRole = role
		tokenInfo.AccountInitPwd = initPwd
		tokenInfo.AccountEmail = email
		tokenInfo.AccountMobile = mobile
		tokenInfo.AccountAvatarUri = avatarUri
		tokenInfo.TokenType = "Bearer"
		tokenInfo.ExpiresAt = expiresAt
		tokenInfo.AccessToken = tokenString
	}
	tokenInfo.RefreshToken.ExpiresAt = expiresAt * a.conf.RefreshExpired
	tokenInfo.RefreshToken.Token = primitive.NewObjectID().Hex()

	return tokenInfo, nil
}

// 公钥解析令牌
func (a *JWTAuth) parseToken(tokenString string) (*claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(tokenString, &claims{},
		func(token *jwt.Token) (interface{}, error) {
			pb, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(a.conf.RSAPublicSecret))
			return pb, nil
		})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

// ParseUserID 解析用户ID
func (a *JWTAuth) ParseUser(ctx context.Context, tokenString string) (string, string, error) {
	if tokenString == "" {
		return "", "", ErrInvalidToken
	}

	claims, err := a.parseToken(tokenString)
	if err != nil {
		return "", "", err
	}
	return claims.UserID, claims.Role, nil
}
