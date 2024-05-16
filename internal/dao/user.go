package dao

import (
	"QA-System/internal/models"

	mysql "QA-System/internal/pkg/database/mysql"
)

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	result := mysql.DB.Model(&models.User{}).Where("username = ?", username).First(&user)
	return &user, result.Error
}

func GetUserByID(id int) (*models.User, error) {
	var user models.User
	result := mysql.DB.Model(&models.User{}).Where("id = ?", id).First(&user)
	return &user, result.Error
}

func CreateUser(user models.User) error {
	result := mysql.DB.Model(models.User{}).Create(&user)
	return result.Error
}