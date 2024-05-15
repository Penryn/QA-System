package admin

import (
	"QA-System/internal/service"
	"QA-System/internal/pkg/code"
	"QA-System/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

type LogData struct {
	Num     int `form:"num" json:"num"`
	LogType int `form:"log_type" binding:"oneof=0 1 2 3 4"`
}

func GetLogMsg(c *gin.Context) {
	var data LogData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		c.Error(&gin.Error{Err: err, Type: gin.ErrorTypeBind})
		utils.JsonErrorResponse(c, code.ParamError)
		return
	}
	//鉴权
	_, err = service.GetUserSession(c)
	if err != nil {
		c.Error(&gin.Error{Err: err, Type: gin.ErrorTypePublic})
		utils.JsonErrorResponse(c, code.NotLogin)
		return
	}
	response, err := service.GetLastLinesFromLogFile("./logs/app.log", data.Num, data.LogType)
	if err != nil {
		c.Error(&gin.Error{Err: err, Type: gin.ErrorTypePublic})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, response)
}
