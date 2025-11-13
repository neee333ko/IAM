package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/authzserver/load"
	"github.com/neee333ko/IAM/pkg/storage"
	"github.com/neee333ko/component-base/pkg/json"
	"github.com/neee333ko/log"
)

func Publish() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := resolvePath(ctx.Request.URL.EscapedPath())

		switch path {
		case "secrets":
			notify(ctx.Request.Method, load.SecretChangedCommand)
		case "policy":
			notify(ctx.Request.Method, load.PolicyChangedCommand)
		}

		ctx.Next()
	}
}

func resolvePath(path string) string {
	parts := strings.Split(path, "/")

	return parts[1]
}

func notify(method string, command load.Command) {
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete {
		client := &storage.RedisCluster{}
		data, _ := json.Marshal(load.Notification{Command: command})
		client.Publish(load.Channel, string(data))

		log.Info("published to redis channel.")
	}
}
