package dao

import (
	"QA-System/internal/models"

	mysql "QA-System/internal/pkg/database/mysql"
)

type Option struct {
	SerialNum int    `json:"serial_num"` //选项序号
	Content   string `json:"content"`    //选项内容
	Img       string `json:"img"`        //图片
}

func CreateOption(option models.Option) error {
	err := mysql.DB.Create(&option).Error
	return err
}

func GetOptionsByQuestionID(questionID int) ([]models.Option, error) {
	var options []models.Option
	err := mysql.DB.Model(models.Option{}).Where("question_id = ?", questionID).Find(&options).Error
	return options, err
}

func DeleteOption(optionID int) error {
	err := mysql.DB.Where("id = ?", optionID).Delete(&models.Option{}).Error
	return err
}