package dao

import (
	"QA-System/internal/models"
	"context"
)

type Option struct {
	SerialNum int    `json:"serial_num"` //选项序号
	Content   string `json:"content"`    //选项内容
	Img       string `json:"img"`        //图片
}

func (d *Dao) CreateOption(ctx context.Context, option models.Option) error {
	err := d.orm.WithContext(ctx).Create(&option).Error
	return err
}

func (d *Dao) GetOptionsByQuestionID(ctx context.Context, questionID int) ([]models.Option, error) {
	var options []models.Option
	err := d.orm.WithContext(ctx).Model(models.Option{}).Where("question_id = ?", questionID).Find(&options).Error
	return options, err
}

func (d *Dao) DeleteOption(ctx context.Context, optionID int) error {
	err := d.orm.WithContext(ctx).Where("id = ?", optionID).Delete(&models.Option{}).Error
	return err
}