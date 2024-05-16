package service

import "QA-System/internal/dao"

func GetAnswerSheetBySurveyID(surveyID int, pageNum int, pageSize int) ([]dao.AnswerSheet, *int64, error) {
	answerSheets, total, err := dao.GetAnswerSheetBySurveyID(surveyID, pageNum, pageSize)
	return answerSheets, total, err
}