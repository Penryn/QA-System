package admin

import (
	"QA-System/internal/dao"
	"QA-System/internal/pkg/code"
	"QA-System/internal/pkg/utils"
	"QA-System/internal/service"
	"errors"
	"math"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// 新建问卷
type CreateSurveyData struct {
	Title     string         `json:"title"`
	Desc      string         `json:"desc" `
	Img       string         `json:"img" `
	Status    int            `json:"status" `
	Time      string         `json:"time"`
	Questions []dao.Question `json:"questions"`
}

func CreateSurvey(c *gin.Context) {
	var data CreateSurveyData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取参数失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ParamError)
		return
	}
	//鉴权
	user, err := service.GetUserSession(c)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取用户缓存信息失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.NotLogin)
		return
	}
	//解析时间转换为中国时间(UTC+8)
	ddlTime, err := time.Parse(time.RFC3339, data.Time)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("时间解析失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	//创建问卷
	err = service.CreateSurvey(user.ID, data.Title, data.Desc, data.Img, data.Questions, data.Status, ddlTime)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("创建问卷失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

// 修改问卷状态
type UpdateSurveyStatusData struct {
	ID     int `json:"id" binding:"required"`
	Status int `json:"status" binding:"required,oneof=1 2"`
}

func UpdateSurveyStatus(c *gin.Context) {
	var data UpdateSurveyStatusData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取参数失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ParamError)
		return
	}
	//鉴权
	user, err := service.GetUserSession(c)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取用户缓存信息失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.NotLogin)
		return
	}
	// 获取问卷
	survey, err := service.GetSurveyByID(data.ID)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取问卷信息失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	//判断权限
	if (user.AdminType != 2) && (user.AdminType != 1 || survey.UserID != user.ID) && !service.UserInManage(user.ID, survey.ID) {
		c.Error(&gin.Error{Err: errors.New("无权限"), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.NoPermission)
		return
	}
	//判断问卷状态
	if survey.Status == data.Status {
		c.Error(&gin.Error{Err: errors.New("问卷状态重复"), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.StatusRepeatError)
		return
	}
	//修改问卷状态
	err = service.UpdateSurveyStatus(data.ID, data.Status)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("修改问卷状态失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

type UpdateSurveyData struct {
	ID        int            `json:"id" binding:"required"`
	Title     string         `json:"title"`
	Desc      string         `json:"desc" `
	Img       string         `json:"img" `
	Time      string         `json:"time"`
	Questions []dao.Question `json:"questions"`
}

func UpdateSurvey(c *gin.Context) {
	var data UpdateSurveyData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取参数失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ParamError)
		return
	}
	//鉴权
	user, err := service.GetUserSession(c)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取用户缓存信息失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.NotLogin)
		return
	}
	// 获取问卷
	survey, err := service.GetSurveyByID(data.ID)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取问卷信息失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	//判断权限
	if (user.AdminType != 2) && (user.AdminType != 1 || survey.UserID != user.ID) && !service.UserInManage(user.ID, survey.ID) {
		c.Error(&gin.Error{Err: errors.New("无权限"), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.NoPermission)
		return
	}
	//判断问卷状态
	if survey.Status != 1 {
		c.Error(&gin.Error{Err: errors.New("问卷状态不为未发布"), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.StatusRepeatError)
		return
	}
	// 判断问卷的填写数量是否为零
	if survey.Num != 0 {
		c.Error(&gin.Error{Err: errors.New("问卷已有填写数量"), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.SurveyNumError)
		return
	}
	//解析时间转换为中国时间(UTC+8)
	ddlTime, err := time.Parse(time.RFC3339, data.Time)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("时间解析失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	//修改问卷
	err = service.UpdateSurvey(data.ID, data.Title, data.Desc, data.Img, data.Questions, ddlTime)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("修改问卷失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

// 删除问卷
type DeleteSurveyData struct {
	ID int `form:"id" binding:"required"`
}

func DeleteSurvey(c *gin.Context) {
	var data DeleteSurveyData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取参数失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ParamError)
		return
	}
	//鉴权
	user, err := service.GetUserSession(c)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取用户缓存信息失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.NotLogin)
		return
	}
	// 获取问卷
	survey, err := service.GetSurveyByID(data.ID)
	if err == gorm.ErrRecordNotFound {
		c.Error(&gin.Error{Err: errors.New("问卷不存在"), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.SurveyNotExist)
		return
	} else if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取问卷信息失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	//判断权限
	if (user.AdminType != 2) && (user.AdminType != 1 || survey.UserID != user.ID) && !service.UserInManage(user.ID, survey.ID) {
		c.Error(&gin.Error{Err: errors.New("无权限"), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.NoPermission)
		return
	}
	//删除问卷
	err = service.DeleteSurvey(data.ID)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("删除问卷失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

// 获取问卷收集数据
type GetSurveyAnswersData struct {
	ID       int `form:"id" binding:"required"`
	PageNum  int `form:"page_num" binding:"required"`
	PageSize int `form:"page_size" binding:"required"`
}

func GetSurveyAnswers(c *gin.Context) {
	var data GetSurveyAnswersData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取参数失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ParamError)
		return
	}
	//鉴权
	user, err := service.GetUserSession(c)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取用户缓存信息失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.NotLogin)
		return
	}
	// 获取问卷
	survey, err := service.GetSurveyByID(data.ID)
	if err == gorm.ErrRecordNotFound {
		c.Error(&gin.Error{Err: errors.New("问卷不存在"), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.SurveyNotExist)
		return
	} else if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取问卷信息失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	//判断权限
	if (user.AdminType != 2) && (user.AdminType != 1 || survey.UserID != user.ID) && !service.UserInManage(user.ID, survey.ID) {
		c.Error(&gin.Error{Err: errors.New("无权限"), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.NoPermission)
		return
	}
	//获取问卷收集数据
	var num *int64
	answers, num, err := service.GetSurveyAnswers(data.ID, data.PageNum, data.PageSize)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取问卷收集数据失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, gin.H{
		"answers_data":   answers,
		"total_page_num": math.Ceil(float64(*num) / float64(data.PageSize)),
	})
}

type GetAllSurveyData struct {
	PageNum  int    `form:"page_num" binding:"required"`
	PageSize int    `form:"page_size" binding:"required"`
	Title    string `form:"title"`
}

func GetAllSurvey(c *gin.Context) {
	var data GetAllSurveyData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取参数失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ParamError)
		return
	}
	user, err := service.GetUserSession(c)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取用户缓存信息失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.NotLogin)
		return
	}
	// 获取问卷
	response := make([]interface{}, 0)
	var totalPageNum *int64
	if user.AdminType == 2 {
		response, totalPageNum, err = service.GetAllSurvey(data.PageNum, data.PageSize, data.Title)
		if err != nil {
			c.Error(&gin.Error{Err: errors.New("获取问卷信息失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
			utils.JsonErrorResponse(c, code.ServerError)
			return
		}
	} else {
		response, err = service.GetAllSurveyByUserID(user.ID)
		if err != nil {
			c.Error(&gin.Error{Err: errors.New("获取问卷信息失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
			utils.JsonErrorResponse(c, code.ServerError)
			return
		}
		managedSurveys, err := service.GetManageredSurveyByUserID(user.ID)
		if err != nil {
			c.Error(&gin.Error{Err: errors.New("获取问卷信息失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
			utils.JsonErrorResponse(c, code.ServerError)
			return
		}
		for _, manage := range managedSurveys {
			managedSurvey, err := service.GetSurveyByID(manage.SurveyID)
			if err != nil {
				c.Error(&gin.Error{Err: errors.New("获取问卷信息失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
				utils.JsonErrorResponse(c, code.ServerError)
				return
			}
			managedSurveyResponse := map[string]interface{}{
				"id":     managedSurvey.ID,
				"title":  managedSurvey.Title,
				"status": managedSurvey.Status,
				"num":    managedSurvey.Num,
			}
			response = append(response, managedSurveyResponse)
		}
		response, totalPageNum = service.ProcessResponse(response, data.PageNum, data.PageSize, data.Title)
	}

	utils.JsonSuccessResponse(c, gin.H{
		"survey_list":    response,
		"total_page_num": math.Ceil(float64(*totalPageNum) / float64(data.PageSize)),
	})
}

type GetSurveyData struct {
	ID int `form:"id" binding:"required"`
}

type SurveyData struct {
	ID        int            `json:"id"`
	Time      string         `json:"time"`
	Desc      string         `json:"desc"`
	Img       string         `json:"img"`
	Questions []dao.Question `json:"questions"`
}

// 管理员获取问卷题面
func GetSurvey(c *gin.Context) {
	var data GetSurveyData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取参数失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ParamError)
		return
	}
	user, err := service.GetUserSession(c)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取用户缓存信息失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.NotLogin)
		return
	}
	// 获取问卷
	survey, err := service.GetSurveyByID(data.ID)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取问卷信息失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	//判断权限
	if (user.AdminType != 2) && (user.AdminType != 1 || survey.UserID != user.ID) && !service.UserInManage(user.ID, survey.ID) {
		c.Error(&gin.Error{Err: errors.New("无权限"), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.NoPermission)
		return
	}
	// 获取相应的问题
	questions, err := service.GetQuestionsBySurveyID(survey.ID)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取问题失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	// 构建问卷响应
	questionsResponse := make([]map[string]interface{}, 0)
	for _, question := range questions {
		options, err := service.GetOptionsByQuestionID(question.ID)
		if err != nil {
			c.Error(&gin.Error{Err: errors.New("获取选项失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
			utils.JsonErrorResponse(c, code.ServerError)
			return
		}
		optionsResponse := make([]map[string]interface{}, 0)
		for _, option := range options {
			optionResponse := map[string]interface{}{
				"img":        option.Img,
				"content":    option.Content,
				"serial_num": option.SerialNum,
			}
			optionsResponse = append(optionsResponse, optionResponse)
		}
		questionMap := map[string]interface{}{
			"id":            question.SerialNum,
			"serial_num":    question.SerialNum,
			"subject":       question.Subject,
			"description":   question.Description,
			"required":      question.Required,
			"unique":        question.Unique,
			"other_option":  question.OtherOption,
			"img":           question.Img,
			"question_type": question.QuestionType,
			"reg":           question.Reg,
			"options":       optionsResponse,
		}
		questionsResponse = append(questionsResponse, questionMap)
	}
	response := map[string]interface{}{
		"id":        survey.ID,
		"title":     survey.Title,
		"time":      survey.Deadline,
		"desc":      survey.Desc,
		"img":       survey.Img,
		"questions": questionsResponse,
	}

	utils.JsonSuccessResponse(c, response)
}

type DownloadFileData struct {
	ID int `form:"id" binding:"required"`
}

// 下载
func DownloadFile(c *gin.Context) {
	var data DownloadFileData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取参数失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ParamError)
		return
	}
	user, err := service.GetUserSession(c)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取用户缓存信息失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.NotLogin)
		return
	}
	// 获取问卷
	survey, err := service.GetSurveyByID(data.ID)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取问卷信息失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	// 判断权限
	if (user.AdminType != 2) && (user.AdminType != 1 || survey.UserID != user.ID) && !service.UserInManage(user.ID, survey.ID) {
		c.Error(&gin.Error{Err: errors.New("无权限"), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.NoPermission)
		return
	}
	// 获取数据
	answers, err := service.GetAllSurveyAnswers(data.ID)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取问卷收集数据失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	url, err := service.HandleDownloadFile(answers, survey)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("文件下载失败" + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, url)
}
