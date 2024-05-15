package admin

import (
	"QA-System/internal/service"
	"QA-System/internal/pkg/utils"
	"QA-System/internal/pkg/code"
	"errors"

	"github.com/gin-gonic/gin"
)

type CreatePermissionData struct {
	UserName string `json:"username"`
	SurveyID int    `json:"survey_id"`
}

func CreatrPermission(c *gin.Context) {
	var data CreatePermissionData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(&gin.Error{Err: err, Type: gin.ErrorTypeBind})
		utils.JsonErrorResponse(c, code.ParamError)
		return
	}
	//鉴权
	admin, err := service.GetUserSession(c)
	if err != nil {
		c.Error(&gin.Error{Err: err, Type: gin.ErrorTypePublic})
		utils.JsonErrorResponse(c, code.NotLogin)
		return
	}
	if admin.AdminType != 2 {
		c.Error(errors.New("没有权限"))
		utils.JsonErrorResponse(c, code.NoPermission)
		return
	}
	user, err := service.GetUserByName(data.UserName)
	if err != nil {
		c.Error(&gin.Error{Err: err, Type: gin.ErrorTypePublic})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	survey, err := service.GetSurveyByID(data.SurveyID)
	if err != nil {
		c.Error(&gin.Error{Err: err, Type: gin.ErrorTypePublic})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	if survey.UserID == user.ID {
		c.Error(errors.New("不能给问卷所有者添加权限"))
		utils.JsonErrorResponse(c, code.PermissionBelong)
		return
	}
	err = service.CheckPermission(user.ID, data.SurveyID)
	if err != nil {
		c.Error(&gin.Error{Err: err, Type: gin.ErrorTypePublic})
		utils.JsonErrorResponse(c, code.PermissionExist)
		return
	}
	//创建权限
	err = service.CreatePermission(user.ID, data.SurveyID)
	if err != nil {
		c.Error(&gin.Error{Err: err, Type: gin.ErrorTypePublic})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

type DeletePermissionData struct {
	UserName string `form:"username"`
	SurveyID int    `form:"survey_id"`
}

func DeletePermission(c *gin.Context) {
	var data DeletePermissionData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		c.Error(&gin.Error{Err: err, Type: gin.ErrorTypePublic})
		utils.JsonErrorResponse(c, code.ParamError)
		return
	}
	//鉴权
	admin, err := service.GetUserSession(c)
	if err != nil {
		c.Error(&gin.Error{Err: err, Type: gin.ErrorTypePublic})
		utils.JsonErrorResponse(c, code.NotLogin)
		return
	}
	if admin.AdminType != 2 {
		c.Error(errors.New("没有权限"))
		utils.JsonErrorResponse(c, code.NoPermission)
		return
	}
	user, err := service.GetUserByName(data.UserName)
	if err != nil {
		c.Error(&gin.Error{Err: err, Type: gin.ErrorTypePublic})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	survey, err := service.GetSurveyByID(data.SurveyID)
	if err != nil {
		c.Error(&gin.Error{Err: err, Type: gin.ErrorTypePublic})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	if survey.UserID == user.ID {
		c.Error(errors.New("不能删除问卷所有者的权限"))
		utils.JsonErrorResponse(c, code.PermissionBelong)
		return
	}
	//删除权限
	err = service.DeletePermission(user.ID, data.SurveyID)
	if err != nil {
		c.Error(&gin.Error{Err: err, Type: gin.ErrorTypePublic})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}
