package dao

import (
	"QA-System/internal/models"
	"context"
)

// GetUserByUsername 根据用户名获取用户
func (d *Dao) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	result := d.orm.WithContext(ctx).Model(&models.User{}).Where("username = ?", username).First(&user)
	return &user, result.Error
}

// GetUserByID 根据用户ID获取用户
func (d *Dao) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	result := d.orm.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).First(&user)
	return &user, result.Error
}

// CreateUser 创建新用户
func (d *Dao) CreateUser(ctx context.Context, user *models.User) error {
	result := d.orm.WithContext(ctx).Model(&models.User{}).Create(user)
	return result.Error
}