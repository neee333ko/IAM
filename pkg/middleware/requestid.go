package middleware

import (
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

const KeyRequestID = "X-Request-ID"

func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rid := ctx.GetHeader(KeyRequestID)

		if rid == "" {
			rid = uuid.NewV4().String()
			ctx.Set(KeyRequestID, rid)
			ctx.Request.Header.Set(KeyRequestID, rid)
		}

		ctx.Header(KeyRequestID, rid)

		ctx.Next()
	}
}

func GetLoggerConfig(formatter gin.LogFormatter, output io.Writer, skippaths []string) gin.LoggerConfig {
	if formatter == nil {
		formatter = LogFormatterWithRequestID
	}

	return gin.LoggerConfig{
		Formatter: formatter,
		Output:    output,
		SkipPaths: skippaths,
	}
}

// return default gin format info plus requestid
func LogFormatterWithRequestID(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}
	return fmt.Sprintf("[GIN] %v | %s |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		param.Request.Header.Get(KeyRequestID),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}

func GetRequestIDFromContext(ctx *gin.Context) string {
	rid, ok := ctx.Get(KeyRequestID)
	if ok {
		if rids, ok := rid.(string); ok {
			return rids
		}
	}

	return ""
}

func GetRequestIDFromHeader(ctx *gin.Context) string {
	return ctx.GetHeader(KeyRequestID)
}
