package service

import (
	"QA-System/internal/dao"
	global "QA-System/internal/global/config"
	"context"
	"time"

	r "QA-System/internal/pkg/redis"
)

var (
	ctx = context.Background()
)

func GetAnswerSheetBySurveyID(surveyID int, pageNum int, pageSize int) ([]dao.AnswerSheet, *int64, error) {
	answerSheets, total, err := dao.GetAnswerSheetBySurveyID(surveyID, pageNum, pageSize)
	return answerSheets, total, err
}

func GetConfigUrl() string {
	url := GetRedis("url")
	if url == "" {
		url=global.Config.GetString("url.host")
		SetRedis("url", url)
	}
	return url
}

func GetConfigKey() string {
	key := GetRedis("key")
	if key == "" {
		key=global.Config.GetString("key")
		SetRedis("key", key)
	}
	return key
}


func SetRedis(key string, value string) bool {
	t := int64(900)
	expire := time.Duration(t) * time.Second
	if err := r.RedisClient.Set(ctx, key, value, expire).Err(); err != nil {
		return false
	}
	return true
}

func GetRedis(key string) string {
	result, err := r.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return result
}