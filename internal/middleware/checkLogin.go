package middlewares

import (
	"QA-System/internal/pkg/code"
	"QA-System/internal/service"
	"QA-System/internal/pkg/utils"
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckLogin(c *gin.Context) {
	isLogin := service.CheckUserSession(c)
	if !isLogin {
		c.Error(errors.New("未登录"))
		utils.JsonErrorResponse(c, code.NotLogin)
		c.Abort()
		return
	}
	c.Next()
}
