package dao

import (
	"QA-System/internal/models"
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockDao struct {
	mock.Mock
}

// user
func (m *MockDao) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockDao) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockDao) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

//answer
func (m *MockDao) SaveAnswerSheet(ctx context.Context, answerSheet AnswerSheet) error {
	args := m.Called(ctx, answerSheet)
	return args.Error(0)
}

func (m *MockDao) GetAnswerSheetBySurveyID(ctx context.Context, surveyID int, pageNum int, pageSize int) ([]AnswerSheet, *int64, error) {
	args := m.Called(ctx, surveyID, pageNum, pageSize)
	return args.Get(0).([]AnswerSheet), args.Get(1).(*int64), args.Error(2)
}

func (m *MockDao) DeleteAnswerSheetBySurveyID(ctx context.Context, surveyID int) error {
	args := m.Called(ctx, surveyID)
	return args.Error(0)
}

//manage
func (m *MockDao) CreateManage(ctx context.Context, id int, surveyID int) error {
	args := m.Called(ctx, id, surveyID)
	return args.Error(0)
}

func (m *MockDao) DeleteManage(ctx context.Context, id int, surveyID int) error {
	args := m.Called(ctx, id, surveyID)
	return args.Error(0)
}

func (m *MockDao) DeleteManageBySurveyID(ctx context.Context, surveyID int) error {
	args := m.Called(ctx, surveyID)
	return args.Error(0)
}

func (m *MockDao) CheckManage(ctx context.Context, id int, surveyID int) error {
	args := m.Called(ctx, id, surveyID)
	return args.Error(0)
}

func (m *MockDao) GetManageByUIDAndSID(ctx context.Context, uid int, sid int) (*models.Manage, error) {
	args := m.Called(ctx, uid, sid)
	return args.Get(0).(*models.Manage), args.Error(1)
}

func (m *MockDao) GetManageByUserID(ctx context.Context, uid int) ([]models.Manage, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).([]models.Manage), args.Error(1)
}

//option
func (m *MockDao) CreateOption(ctx context.Context, option *models.Option) error {
	args := m.Called(ctx, option)
	return args.Error(0)
}

func (m *MockDao) GetOptionsByQuestionID(ctx context.Context, questionID int) ([]models.Option, error) {
	args := m.Called(ctx, questionID)
	return args.Get(0).([]models.Option), args.Error(1)
}

func (m *MockDao) DeleteOption(ctx context.Context, optionID int) error {
	args := m.Called(ctx, optionID)
	return args.Error(0)
}

//question
func (m *MockDao) CreateQuestion(ctx context.Context, question *models.Question) error {
	args := m.Called(ctx, question)
	return args.Error(0)
}

func (m *MockDao) GetQuestionsBySurveyID(ctx context.Context, surveyID int) ([]models.Question, error) {
	args := m.Called(ctx, surveyID)
	return args.Get(0).([]models.Question), args.Error(1)
}

func (m *MockDao) GetQuestionByID(ctx context.Context, questionID int) (*models.Question, error) {
	args := m.Called(ctx, questionID)
	return args.Get(0).(*models.Question), args.Error(1)
}

func (m *MockDao) DeleteQuestion(ctx context.Context, questionID int) error {
	args := m.Called(ctx, questionID)
	return args.Error(0)
}

func (m *MockDao) DeleteQuestionBySurveyID(ctx context.Context, surveyID int) error {
	args := m.Called(ctx, surveyID)
	return args.Error(0)
}

//survey
func (m *MockDao) CreateSurvey(ctx context.Context, survey *models.Survey) error {
	args := m.Called(ctx, survey)
	return args.Error(0)
}

func (m *MockDao) GetSurveyByID(ctx context.Context, surveyID int) (*models.Survey, error) {
	args := m.Called(ctx, surveyID)
	return args.Get(0).(*models.Survey), args.Error(1)
}

func (m *MockDao) GetSurveyByTitle(ctx context.Context, title string, num, size int) ([]models.Survey, *int64, error) {
	args := m.Called(ctx, title, num, size)
	return args.Get(0).([]models.Survey), args.Get(1).(*int64), args.Error(2)
}

func (m *MockDao) DeleteSurvey(ctx context.Context, surveyID int) error {
	args := m.Called(ctx, surveyID)
	return args.Error(0)
}

func (m *MockDao) GetAllSurveyByUserID(ctx context.Context, uid int) ([]models.Survey, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).([]models.Survey), args.Error(1)
}

func (m *MockDao) IncreaseSurveyNum(ctx context.Context, sid int) error {
	args := m.Called(ctx, sid)
	return args.Error(0)
}

func (m *MockDao) UpdateSurveyStatus(ctx context.Context, surveyID int, status int) error {
	args := m.Called(ctx, surveyID, status)
	return args.Error(0)
}

func (m *MockDao) UpdateSurvey(ctx context.Context, id int, title, desc, img string, deadline time.Time) error {
	args := m.Called(ctx, id, title, desc, img, deadline)
	return args.Error(0)
}
