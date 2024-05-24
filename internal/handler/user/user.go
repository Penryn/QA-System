package user

import (
	"QA-System/internal/dao"
	"QA-System/internal/handler/queue"
	"QA-System/internal/pkg/code"
	"QA-System/internal/pkg/queue/asynq"
	"QA-System/internal/pkg/utils"
	"QA-System/internal/service"
	"errors"

	"time"

	"github.com/gin-gonic/gin"
)

type SubmitServeyData struct {
	ID            int                 `json:"id" binding:"required"`
	QuestionsList []dao.QuestionsList `json:"questions_list"`
}

func SubmitSurvey(c *gin.Context) {
	var data SubmitServeyData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取参数失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ParamError)
		return
	}
	// 判断问卷问题和答卷问题数目是否一致
	survey, err := service.GetSurveyByID(data.ID)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取问卷失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	questions, err := service.GetQuestionsBySurveyID(survey.ID)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取问题失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	if len(questions) != len(data.QuestionsList) {
		c.Error(&gin.Error{Err: errors.New("问题数量不一致"), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	// 判断填写时间是否在问卷有效期内
	if !survey.Deadline.IsZero() && survey.Deadline.Before(time.Now()) {
		c.Error(&gin.Error{Err: errors.New("填写时间已过"), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.TimeBeyondError)
		return
	}
	// 逐个判断问题答案
	for _, q := range data.QuestionsList {
		question, err := service.GetQuestionByID(q.QuestionID)
		if err != nil {
			c.Error(&gin.Error{Err: errors.New("获取问题失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
			utils.JsonErrorResponse(c, code.ServerError)
			return
		}
		if question.SerialNum != q.SerialNum {
			c.Error(&gin.Error{Err: errors.New("问题序号不一致"), Type: gin.ErrorTypeAny})
			utils.JsonErrorResponse(c, code.ServerError)
			return
		}
		if question.SurveyID != survey.ID {
			c.Error(&gin.Error{Err: errors.New("问题不属于该问卷"), Type: gin.ErrorTypeAny})
			utils.JsonErrorResponse(c, code.ServerError)
			return
		}
		// 判断必填字段是否为空
		if question.Required && q.Answer == "" {
			c.Error(&gin.Error{Err: errors.New("必填字段为空"), Type: gin.ErrorTypeAny})
			utils.JsonErrorResponse(c, code.ServerError)
			return
		}
		// 判断唯一字段是否唯一
		if question.Unique {
			unique, err := service.CheckUnique(data.ID, q.QuestionID, question.SerialNum, q.Answer)
			if err != nil {
				c.Error(&gin.Error{Err: errors.New("唯一字段检查失败"), Type: gin.ErrorTypeAny})
				utils.JsonErrorResponse(c, code.ServerError)
				return
			}
			if !unique {
				c.Error(&gin.Error{Err: errors.New("唯一字段不唯一"), Type: gin.ErrorTypeAny})
				utils.JsonErrorResponse(c, code.UniqueError)
				return
			}

		}
	}
	// 创建并入队任务
	task, err := queue.NewSubmitSurveyTask(data.ID, data.QuestionsList)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("创建任务失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}

	_, err = asynq.Client.Enqueue(task)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("任务入队失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, nil)
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

// 用户获取问卷
func GetSurvey(c *gin.Context) {
	var data GetSurveyData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取参数失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ParamError)
		return
	}
	// 获取问卷
	survey, err := service.GetSurveyByID(data.ID)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取问卷失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	// 判断填写时间是否在问卷有效期内
	if !survey.Deadline.IsZero() && survey.Deadline.Before(time.Now()) {
		c.Error(&gin.Error{Err: errors.New("填写时间已过"), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.TimeBeyondError)
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
			"id":            question.ID,
			"serial_num":    question.SerialNum,
			"subject":       question.Subject,
			"describe":      question.Description,
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

// 上传图片
func UploadImg(c *gin.Context) {
	url, err := service.HandleImgUpload(c)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("上传图片失败" + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, url)
}
