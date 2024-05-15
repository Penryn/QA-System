package middlewares

import (
	"QA-System/internal/pkg/code"
	"QA-System/internal/pkg/log"

	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Handle panic and log the error
				stack := debug.Stack()
				log.Logger.Error("Panic recovered",
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.Any("panic", r),
					zap.ByteString("stacktrace", stack),
				)

				c.JSON(http.StatusInternalServerError, gin.H{
					"code":  http.StatusInternalServerError,
					"msg": code.ServerError.Msg,
				})
				c.Abort()
			}
		}()

		c.Next()
		if length := len(c.Errors); length > 0 {
			e := c.Errors[length-1]
			err := e.Err
			if err != nil {
				// TODO 建立日志系统
				var logLevel zapcore.Level
				switch e.Type {
				case gin.ErrorTypePublic:
					logLevel = zapcore.ErrorLevel
				case gin.ErrorTypeBind:
					logLevel = zapcore.WarnLevel
				case gin.ErrorTypePrivate:
					logLevel = zapcore.DebugLevel
				default:
					logLevel = zapcore.InfoLevel
				}
				log.Logger.Check(logLevel, "Error reported").Write(
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.Error(err),
				)
				return
			}
		}
	}
}


// HandleNotFound
//
//	404处理
func HandleNotFound(c *gin.Context) {
	err := code.NotFound
	c.JSON(err.StatusCode, err)
}
