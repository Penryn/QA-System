package dao

import (
	"QA-System/internal/models"
	"time"

	mysql "QA-System/internal/pkg/database/mysql"

	"gorm.io/gorm"
)

func CreateSurvey(survey models.Survey) error {
	err := mysql.DB.Create(&survey).Error
	return err
}

func UpdateSurveyStatus(surveyID int, status int) error {
	err := mysql.DB.Model(&models.Survey{}).Where("id = ?", surveyID).Update("status", status).Error
	return err
}

func UpdateSurvey(id int, title, desc, img string, deadline time.Time) error {
	err := mysql.DB.Model(&models.Survey{}).Where("id = ?", id).Updates(models.Survey{Title: title, Desc: desc, Img: img, Deadline: deadline}).Error
	return err
}

func GetAllSurveyByUserID(userId int) ([]models.Survey, error) {
	var surveys []models.Survey
	err := mysql.DB.Model(models.Survey{}).Where("user_id = ?", userId).
		Order("CASE WHEN status = 2 THEN 0 ELSE 1 END, id DESC").Find(&surveys).Error
	return surveys, err
}

func GetSurveyByID(surveyID int) (*models.Survey, error) {
	var survey models.Survey
	err := mysql.DB.Where("id = ?", surveyID).First(&survey).Error
	return &survey, err
}

func GetSurveyByTitle(title string,num,size int) ([]models.Survey,*int64, error) {
	var surveys []models.Survey
	var sum int64
	err := mysql.DB.Model(models.Survey{}).Where("title like ?", "%"+title+"%").Order("CASE WHEN status = 2 THEN 0 ELSE 1 END, id DESC").Count(&sum).Limit(size).Offset((num-1)*size).Find(&surveys).Error
	return surveys, &sum, err
}

func IncreaseSurveyNum(sid int) error {
	err :=mysql.DB.Model(&models.Survey{}).Where("id = ?", sid).Update("num", gorm.Expr("num + ?", 1)).Error
	return err
}

func DeleteSurvey(surveyID int) error {
	err := mysql.DB.Where("id = ?", surveyID).Delete(&models.Survey{}).Error
	return err
}


