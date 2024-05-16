package dao


import (
	"QA-System/internal/models"

	mysql "QA-System/internal/pkg/database/mysql"
)


type Question struct {
	ID           int      `json:"id"`
	SerialNum    int      `json:"serial_num"`    //题目序号
	Subject      string   `json:"subject"`       //问题
	Description  string   `json:"description"`   //问题描述
	Img          string   `json:"img"`           //图片
	Required     bool     `json:"required"`      //是否必填
	Unique       bool     `json:"unique"`        //是否唯一
	OtherOption  bool     `json:"other_option"`  //是否有其他选项
	QuestionType int      `json:"question_type"` //问题类型 1单选2多选3填空4简答5图片
	Reg          string   `json:"reg"`           //正则表达式
	Options      []Option `json:"options"`       //选项
}

type QuestionsList struct {
	QuestionID int    `json:"question_id" binding:"required"`
	SerialNum  int    `json:"serial_num"`
	Answer     string `json:"answer"`
}

func CreateQuestion(question models.Question) error {
	err := mysql.DB.Create(&question).Error
	return err
}

func GetQuestionsBySurveyID(surveyID int) ([]models.Question, error) {
	var questions []models.Question
	err := mysql.DB.Model(models.Question{}).Where("survey_id = ?", surveyID).Find(&questions).Error
	return questions, err

}

func GetQuestionByID(questionID int) (*models.Question, error) {
	var question models.Question
	err := mysql.DB.Where("id = ?", questionID).First(&question).Error
	return &question, err
}

func DeleteQuestion(questionID int) error {
	err := mysql.DB.Where("id = ?", questionID).Delete(&models.Question{}).Error
	return err
}

func DeleteQuestionBySurveyID(surveyID int) error {
	err := mysql.DB.Where("survey_id = ?", surveyID).Delete(&models.Question{}).Error
	return err
}

