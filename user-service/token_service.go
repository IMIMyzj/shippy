package main

import (
	pb "shippy/user-service/proto/user"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Authable interface {
	Decode(tokenStr string) (*CustomClaims, error)
	Encode(user *pb.User) (string, error)
}

// 定义加盐哈希密码时所用的盐  要保证其生成和保存都足够安全  比如使用 md5 来生成
var privateKey = []byte("`xs#a_1-!")

// 自定义的metadata, 在加密后变成payload作为JWT的第二部分返回给客户端
type CustomClaims struct {
	User *pb.User
	// 使用标准的 payload
	jwt.StandardClaims
}

type TokenService struct {
	repo Repository
}

// 将 JWT 字符串解密为 CustomClaims 对象
func (srv *TokenService) Decode(tokenStr string) (*CustomClaims, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return privateKey, nil
	})
	// 解密转换类型并返回，t.Claims是JWT中间部分的字符串
	if claims, ok := t.Claims.(*CustomClaims); ok && t.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

// 将 User 用户信息加密为 JWT 字符串
func (srv *TokenService) Encode(user *pb.User) (string, error) {
	// 三天后过期
	expireTime := time.Now().Add(time.Hour * 24 * 3).Unix()
	claims := CustomClaims{
		user,
		jwt.StandardClaims{
			Issuer:    "go.micro.srv.user", // 签发者
			ExpiresAt: expireTime,          // 给定过期时间
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // 使用HS256算法加密
	return jwtToken.SignedString(privateKey)                      // 在后面加入私钥
}
