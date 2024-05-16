package dao

import (
	"QA-System/internal/models"

	mysql "QA-System/internal/pkg/database/mysql"
)


func CreateManage(id int, surveyID int) error {
	err := mysql.DB.Create(&models.Manage{UserID: id, SurveyID: surveyID}).Error
	return err
}

func DeleteManage(id int, surveyID int) error {
	err := mysql.DB.Where("user_id = ? AND survey_id = ?", id, surveyID).Delete(&models.Manage{}).Error
	return err
}

func DeleteManageBySurveyID(surveyID int) error {
	err := mysql.DB.Where("survey_id = ?", surveyID).Delete(&models.Manage{}).Error
	return err
}

func CheckManage(id int, surveyID int) error {
	var manage models.Manage
	err := mysql.DB.Where("user_id = ? AND survey_id = ?", id, surveyID).First(&manage).Error
	return err
}

func GetManageByUIDAndSID(uid int, sid int) (*models.Manage, error) {
	var manage models.Manage
	err := mysql.DB.Where("user_id = ? AND survey_id = ?", uid, sid).First(&manage).Error
	return &manage, err
}

func GetManageByUserID(uid int) ([]models.Manage, error) {
	var manages []models.Manage
	err := mysql.DB.Where("user_id = ?", uid).Find(&manages).Error
	return manages, err
}
