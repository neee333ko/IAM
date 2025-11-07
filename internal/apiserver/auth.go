package apiserver

import (
	"context"
	"time"

	ginjwt "github.com/appleboy/gin-jwt/v3"
	jwtcore "github.com/appleboy/gin-jwt/v3/core"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/neee333ko/IAM/internal/apiserver/store"
	"github.com/neee333ko/IAM/internal/pkg/code"
	"github.com/neee333ko/IAM/internal/pkg/middleware"
	"github.com/neee333ko/IAM/internal/pkg/middleware/auth"
	"github.com/neee333ko/component-base/pkg/core"
	v1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/errors"
	"github.com/neee333ko/log"
	"github.com/spf13/viper"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	JWTAuth   *auth.JwtAuth
	BasicAuth *auth.BasicAuth
	AutoAuth  *auth.AutoAuth
)

func InitAuth() {
	m := &ginjwt.GinJWTMiddleware{
		Realm:            viper.GetString("jwt.realm"),
		SigningAlgorithm: viper.GetString("jwt.alg"),
		Key:              []byte(viper.GetString("jwt.key")),
		Timeout:          time.Duration(viper.GetInt64("jwt.timeout")),
		MaxRefresh:       time.Duration(viper.GetInt64("jwt.maxRefresh")),
		Authenticator:    authenticator(),
		Unauthorized:     unauthorized(),
		PayloadFunc:      payloadFunc(),
		LoginResponse:    loginResponse(),
		LogoutResponse:   logoutResponse(),
		RefreshResponse:  loginResponse(),
		IdentityKey:      middleware.KeyUsername,
	}

	JWTAuth = auth.NewJwtAuth(m)

	client := store.Client()

	BasicAuth = auth.NewBasicAuth(Compare(client))

	AutoAuth = auth.NewAutoAuth(BasicAuth, JWTAuth)
}

func authenticator() func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		var user *User = new(User)
		if err := c.ShouldBindJSON(user); err != nil {
			return nil, err
		}

		client := store.Client()

		u, err := client.NewUserStore().Get(c, user.Username, &v1.GetOptions{})
		if err != nil {
			return nil, err
		}

		if err := u.Compare(user.Password); err != nil {
			return nil, err
		}

		return user, nil
	}
}

func unauthorized() func(c *gin.Context, _ int, message string) {
	return func(c *gin.Context, _ int, message string) {
		core.WriteResponse(c, errors.WithCode(code.ErrUnknown, message), nil)
	}
}

func payloadFunc() func(data any) jwt.MapClaims {
	return func(data any) jwt.MapClaims {
		if v, ok := data.(*User); ok {
			return jwt.MapClaims{
				middleware.KeyUsername: v.Username,
			}
		}
		return jwt.MapClaims{}
	}
}

func loginResponse() func(c *gin.Context, token *jwtcore.Token) {
	return func(c *gin.Context, token *jwtcore.Token) {
		core.WriteResponse(c, nil, token)
	}
}

func logoutResponse() func(c *gin.Context) {
	return func(c *gin.Context) {
		core.WriteResponse(c, nil, nil)
	}
}

func Compare(s store.Factory) func(ctx context.Context, name string, password string) bool {
	return func(ctx context.Context, name, password string) bool {
		user, err := s.NewUserStore().Get(ctx, name, &v1.GetOptions{})
		if err != nil {
			log.Errorf("error when getting user:%s from db\n", name)
			return false
		}

		if err := user.Compare(password); err != nil {
			return false
		}

		return true
	}
}
