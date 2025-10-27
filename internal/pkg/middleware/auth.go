package middleware

import "github.com/gin-gonic/gin"

type Auth interface {
	AuthFunc() gin.HandlerFunc
}

type AuthOperator struct {
	a Auth
}

func (op *AuthOperator) SetAuth(a Auth) {
	op.a = a
}

func (op *AuthOperator) Operate() gin.HandlerFunc {
	return op.a.AuthFunc()
}
