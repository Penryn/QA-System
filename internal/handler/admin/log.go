package admin

import (
	"QA-System/internal/pkg/code"
	"QA-System/internal/pkg/utils"
	"QA-System/internal/service"
	"errors"

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
		c.Error(&gin.Error{Err: errors.New("获取参数失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ParamError)
		return
	}
	//鉴权
	_, err = service.GetUserSession(c)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取用户缓存信息失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.NotLogin)
		return
	}
	response, err := service.GetLastLinesFromLogFile(data.Num, data.LogType)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取日志失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}

	for i := 0; i < len(response)/2; i++ {
		j := len(response) - i - 1
		response[i], response[j] = response[j], response[i]
	}
	utils.JsonSuccessResponse(c, response)
}
