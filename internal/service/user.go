package service

import (
	"QA-System/internal/dao"
	"QA-System/internal/models"
	"QA-System/internal/pkg/log"
	"time"



)

func GetSurveyByID(id int) (*models.Survey, error) {
	survey, err := dao.GetSurveyByID(id)
	return survey, err
}

func GetQuestionsBySurveyID(id int) ([]models.Question, error) {
	var questions []models.Question
	questions, err := dao.GetQuestionsBySurveyID(id)
	return questions, err
}

func GetOptionsByQuestionID(questionId int) ([]models.Option, error) {
	options,err:=dao.GetOptionsByQuestionID(questionId)
	return options, err
}

func GetQuestionByID(id int) (*models.Question, error) {
	question, err := dao.GetQuestionByID(id)
	return question, err
}

func CheckUnique(sid int, qid int, serial_num int, content string) (bool, error) {
	var answerSheets []dao.AnswerSheet
	answerSheets, _, err := GetAnswerSheetBySurveyID(sid, 0, 0)
	if err != nil {
		return false, err
	}

	for _, answerSheet := range answerSheets {
		for _, answer := range answerSheet.Answers {
			if answer.QuestionID == qid && answer.SerialNum == serial_num && answer.Content == content {
				return false, nil
			}
		}
	}
	return true, nil
}

func SubmitSurvey(sid int, data []dao.QuestionsList) error {
	var answerSheet dao.AnswerSheet
	answerSheet.SurveyID = sid
	answerSheet.Time = time.Now().Format("2006-01-02 15:04:05")
	for _, q := range data {
		var answer dao.Answer
		answer.QuestionID = q.QuestionID
		answer.SerialNum = q.SerialNum
		answer.Content = q.Answer
		answerSheet.Answers = append(answerSheet.Answers, answer)
	}
	err := SaveAnswerSheet(answerSheet)
	if err != nil {
		return err
	}
	err = dao.IncreaseSurveyNum(sid)
	if err != nil {
		return err
	}
	return nil
}

func SaveAnswerSheet(answerSheet dao.AnswerSheet) error {
	err := dao.SaveAnswerSheet(answerSheet)
	if err != nil {
		log.Logger.Error(err.Error())
	}
	return nil
}