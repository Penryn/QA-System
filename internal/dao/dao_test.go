package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"QA-System/internal/models"
	mongodb "QA-System/internal/pkg/database/mongodb"
	mysql "QA-System/internal/pkg/database/mysql"
	"QA-System/internal/pkg/log"

	"github.com/smartystreets/goconvey/convey"
)

//// mock单元测试
func TestDao(t *testing.T) {
	convey.Convey("Given a MockDao", t, func() {
		mockDao := new(MockDao)

		convey.Convey("When CreateUser is called", func() {
			user := &models.User{Username: "testuser", Password: "password"}
			mockDao.On("CreateUser", context.Background(), user).Return(nil)

			err := mockDao.CreateUser(context.Background(), user)

			convey.Convey("Then the user should be created successfully", func() {
				convey.So(err, convey.ShouldBeNil)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When GetUserByUsername is called", func() {
			user := &models.User{Username: "testuser"}
			mockDao.On("GetUserByUsername", context.Background(), user.Username).Return(user, nil)
			result, err := mockDao.GetUserByUsername(context.Background(), "testuser")

			convey.Convey("Then the correct user should be returned", func() {
				convey.So(err, convey.ShouldBeNil)
				convey.So(result.Username, convey.ShouldEqual, "testuser")
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When GetUserByID is called", func() {
			user := &models.User{ID: 1}
			mockDao.On("GetUserByID", context.Background(), user.ID).Return(user, nil)
			result, err := mockDao.GetUserByID(context.Background(), 1)

			convey.Convey("Then the correct user should be returned", func() {
				convey.So(err, convey.ShouldBeNil)
				convey.So(result.ID, convey.ShouldEqual, 1)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When SaveAnswerSheet is called", func() {
			answerSheet := AnswerSheet{SurveyID: 1, Time: "2021-01-01 00:00:00"}
			mockDao.On("SaveAnswerSheet", context.Background(), answerSheet).Return(nil)

			err := mockDao.SaveAnswerSheet(context.Background(), answerSheet)

			convey.Convey("Then the answer sheet should be saved successfully", func() {
				convey.So(err, convey.ShouldBeNil)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When GetAnswerSheetBySurveyID is called", func() {
			answerSheet := AnswerSheet{SurveyID: 1, Time: "2021-01-01 00:00:00"}
			mockDao.On("GetAnswerSheetBySurveyID", context.Background(), answerSheet.SurveyID, 1, 1).Return([]AnswerSheet{answerSheet}, new(int64), nil)

			result, _, err := mockDao.GetAnswerSheetBySurveyID(context.Background(), answerSheet.SurveyID, 1, 1)

			convey.Convey("Then the correct answer sheet should be returned", func() {
				convey.So(err, convey.ShouldBeNil)
				convey.So(result[0].SurveyID, convey.ShouldEqual, 1)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When DeleteAnswerSheetBySurveyID is called", func() {
			servey:=models.Survey{ID:1}
			mockDao.On("DeleteAnswerSheetBySurveyID", context.Background(), servey.ID).Return(nil)

			err := mockDao.DeleteAnswerSheetBySurveyID(context.Background(), servey.ID)

			convey.Convey("Then the answer sheet should be deleted successfully", func() {
				convey.So(err, convey.ShouldBeNil)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When CreateManage is called", func() {
			mockDao.On("CreateManage", context.Background(), 1, 1).Return(nil)

			err := mockDao.CreateManage(context.Background(), 1, 1)

			convey.Convey("Then the manage should be created successfully", func() {
				convey.So(err, convey.ShouldBeNil)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When DeleteManage is called", func() {
			mockDao.On("DeleteManage", context.Background(), 1, 1).Return(nil)

			err := mockDao.DeleteManage(context.Background(), 1, 1)

			convey.Convey("Then the manage should be deleted successfully", func() {
				convey.So(err, convey.ShouldBeNil)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When DeleteManageBySurveyID is called", func() {
			servey:=models.Survey{ID:1}
			mockDao.On("DeleteManageBySurveyID", context.Background(), servey.ID).Return(nil)

			err := mockDao.DeleteManageBySurveyID(context.Background(), servey.ID)

			convey.Convey("Then the manage should be deleted successfully", func() {
				convey.So(err, convey.ShouldBeNil)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When CheckManage is called", func() {
			mockDao.On("CheckManage", context.Background(), 1, 1).Return(nil)

			err := mockDao.CheckManage(context.Background(), 1, 1)

			convey.Convey("Then the manage should be checked successfully", func() {
				convey.So(err, convey.ShouldBeNil)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When GetManageByUIDAndSID is called", func() {
			manage := &models.Manage{ID: 1, UserID: 1, SurveyID: 1}
			mockDao.On("GetManageByUIDAndSID", context.Background(), manage.UserID, manage.SurveyID).Return(manage, nil)

			result, err := mockDao.GetManageByUIDAndSID(context.Background(),manage.UserID, manage.SurveyID)

			convey.Convey("Then the correct manage should be returned", func() {
				convey.So(err, convey.ShouldBeNil)
				convey.So(result.ID, convey.ShouldEqual, 1)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When GetManageByUserID is called", func() {
			manage := models.Manage{ID: 1, UserID: 1, SurveyID: 1}
			mockDao.On("GetManageByUserID", context.Background(), manage.UserID).Return([]models.Manage{manage}, nil)

			result, err := mockDao.GetManageByUserID(context.Background(), manage.UserID)

			convey.Convey("Then the correct manage should be returned", func() {
				convey.So(err, convey.ShouldBeNil)
				convey.So(result[0].ID, convey.ShouldEqual, 1)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When CreateOption is called", func() {
			option := &models.Option{QuestionID: 1, Content: "content"}
			mockDao.On("CreateOption", context.Background(), option).Return(nil)

			err := mockDao.CreateOption(context.Background(), option)

			convey.Convey("Then the option should be created successfully", func() {
				convey.So(err, convey.ShouldBeNil)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When GetOptionsByQuestionID is called", func() {
			option := models.Option{ID: 1, QuestionID: 1, Content: "content"}
			mockDao.On("GetOptionsByQuestionID", context.Background(), option.QuestionID).Return([]models.Option{option}, nil)

			result, err := mockDao.GetOptionsByQuestionID(context.Background(), option.QuestionID)

			convey.Convey("Then the correct option should be returned", func() {
				convey.So(err, convey.ShouldBeNil)
				convey.So(result[0].ID, convey.ShouldEqual, 1)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When DeleteOption is called", func() {
			mockDao.On("DeleteOption", context.Background(), 1).Return(nil)

			err := mockDao.DeleteOption(context.Background(), 1)

			convey.Convey("Then the option should be deleted successfully", func() {
				convey.So(err, convey.ShouldBeNil)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When CreateQuestion is called", func() {
			question := &models.Question{SurveyID: 1, Subject: "subject"}
			mockDao.On("CreateQuestion", context.Background(), question).Return(nil)

			err := mockDao.CreateQuestion(context.Background(), question)

			convey.Convey("Then the question should be created successfully", func() {
				convey.So(err, convey.ShouldBeNil)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When GetQuestionsBySurveyID is called", func() {
			question := models.Question{ID: 1, SurveyID: 1}
			mockDao.On("GetQuestionsBySurveyID", context.Background(), question.SurveyID).Return([]models.Question{question}, nil)
		
			result, err := mockDao.GetQuestionsBySurveyID(context.Background(), question.SurveyID)
		
			convey.Convey("Then the correct question should be returned", func() {
				convey.So(err, convey.ShouldBeNil)
				convey.So(len(result), convey.ShouldBeGreaterThan, 0)
				convey.So(result[0].ID, convey.ShouldEqual, 1)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When GetQuestionByID is called", func() {
			question := &models.Question{ID: 1, SurveyID: 1, Subject: "subject"}
			mockDao.On("GetQuestionByID", context.Background(), question.ID).Return(question, nil)

			result, err := mockDao.GetQuestionByID(context.Background(), question.ID)

			convey.Convey("Then the correct question should be returned", func() {
				convey.So(err, convey.ShouldBeNil)
				convey.So(result.ID, convey.ShouldEqual, 1)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When DeleteQuestion is called", func() {
			mockDao.On("DeleteQuestion", context.Background(), 1).Return(nil)

			err := mockDao.DeleteQuestion(context.Background(), 1)

			convey.Convey("Then the question should be deleted successfully", func() {
				convey.So(err, convey.ShouldBeNil)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When DeleteQuestionBySurveyID is called", func() {
			servey:=models.Survey{ID:1}
			mockDao.On("DeleteQuestionBySurveyID", context.Background(), servey.ID).Return(nil)

			err := mockDao.DeleteQuestionBySurveyID(context.Background(), servey.ID)

			convey.Convey("Then the question should be deleted successfully", func() {
				convey.So(err, convey.ShouldBeNil)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When CreateSurvey is called", func() {
			survey := &models.Survey{Title: "title"}
			mockDao.On("CreateSurvey", context.Background(), survey).Return(nil)

			err := mockDao.CreateSurvey(context.Background(), survey)

			convey.Convey("Then the survey should be created successfully", func() {
				convey.So(err, convey.ShouldBeNil)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When GetSurveyByID is called", func() {
			survey := &models.Survey{ID: 1, Title: "title"}
			mockDao.On("GetSurveyByID", context.Background(), survey.ID).Return(survey, nil)

			result, err := mockDao.GetSurveyByID(context.Background(), survey.ID)

			convey.Convey("Then the correct survey should be returned", func() {
				convey.So(err, convey.ShouldBeNil)
				convey.So(result.ID, convey.ShouldEqual, 1)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When GetSurveyByTitle is called", func() {
			survey := models.Survey{ID: 1, Title: "title"}
			mockDao.On("GetSurveyByTitle", context.Background(), "title", 1, 1).Return([]models.Survey{survey}, new(int64), nil)

			result, _, err := mockDao.GetSurveyByTitle(context.Background(), "title", 1, 1)

			convey.Convey("Then the correct survey should be returned", func() {
				convey.So(err, convey.ShouldBeNil)
				convey.So(result[0].ID, convey.ShouldEqual, 1)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When DeleteSurvey is called", func() {
			mockDao.On("DeleteSurvey", context.Background(), 1).Return(nil)

			err := mockDao.DeleteSurvey(context.Background(), 1)

			convey.Convey("Then the survey should be deleted successfully", func() {
				convey.So(err, convey.ShouldBeNil)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When IncreaseSurveyNum is called", func() {
			mockDao.On("IncreaseSurveyNum", context.Background(), 1).Return(nil)

			err := mockDao.IncreaseSurveyNum(context.Background(), 1)

			convey.Convey("Then the survey num should be increased successfully", func() {
				convey.So(err, convey.ShouldBeNil)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When UpdateSurveyStatus is called", func() {
			mockDao.On("UpdateSurveyStatus", context.Background(), 1, 1).Return(nil)

			err := mockDao.UpdateSurveyStatus(context.Background(), 1, 1)

			convey.Convey("Then the survey status should be updated successfully", func() {
				convey.So(err, convey.ShouldBeNil)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When UpdateSurvey is called", func() {
			mockDao.On("UpdateSurvey", context.Background(), 1, "title", "desc", "img", time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)).Return(nil)

			err := mockDao.UpdateSurvey(context.Background(), 1, "title", "desc", "img", time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))

			convey.Convey("Then the survey should be updated successfully", func() {
				convey.So(err, convey.ShouldBeNil)
				mockDao.AssertExpectations(t)
			})
		})

		convey.Convey("When GetAllSurveyByUserID is called", func() {
			survey := models.Survey{ID: 1, Title: "title"}
			mockDao.On("GetAllSurveyByUserID", context.Background(), 1).Return([]models.Survey{survey}, nil)

			result, err := mockDao.GetAllSurveyByUserID(context.Background(), 1)

			convey.Convey("Then the correct survey should be returned", func() {
				convey.So(err, convey.ShouldBeNil)
				convey.So(result[0].ID, convey.ShouldEqual, 1)
				mockDao.AssertExpectations(t)
			})
		})

	})
}

//// 性能测试
// mysql 增加
func BenchmarkCreateSurvey(b *testing.B) {
	log.ZapInit()
	d := New(mysql.MysqlInit(), mongodb.MongodbInit())

	b.ResetTimer()
	title :=fmt.Sprintf("title%d",time.Now().Unix())
	survey := models.Survey{Title: title, Desc: "desc", Img: "img", Deadline: time.Now()}
	for i := 0; i < b.N; i++ {
		d.CreateSurvey(context.Background(), survey)
	}
}

// mysql 查询
func BenchmarkGetSurveyByTitle(b *testing.B) {
	log.ZapInit()
	d := New(mysql.MysqlInit(), mongodb.MongodbInit())

	b.ResetTimer()
	title :=fmt.Sprintf("title%d",time.Now().Unix())
	for i := 0; i < b.N; i++ {
		d.GetSurveyByTitle(context.Background(), title, 1, 1)
	}
}

// mysql 删除
func BenchmarkDeleteSurvey(b *testing.B) {
	log.ZapInit()
	d := New(mysql.MysqlInit(), mongodb.MongodbInit())

	b.ResetTimer()
	for i := 15; i < b.N; i++ {
		d.DeleteSurvey(context.Background(), i)
	}
}

// mysql 更新
func BenchmarkUpdateSurvey(b *testing.B) {
	log.ZapInit()
	d := New(mysql.MysqlInit(), mongodb.MongodbInit())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.UpdateSurvey(context.Background(), 1, "title", "desc", "img", time.Now())
	}
}

// mongodb 增加
func BenchmarkSaveAnswerSheet(b *testing.B) {
	log.ZapInit()
	d := New(mysql.MysqlInit(), mongodb.MongodbInit())

	b.ResetTimer()
	answerSheet := AnswerSheet{SurveyID: 1, Time: "2021-01-01 00:00:00"}
	for i := 0; i < b.N; i++ {
		d.SaveAnswerSheet(context.Background(), answerSheet)
	}
}

// mongodb 查询
func BenchmarkGetAnswerSheetBySurveyID(b *testing.B) {
	log.ZapInit()
	d := New(mysql.MysqlInit(), mongodb.MongodbInit())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.GetAnswerSheetBySurveyID(context.Background(), 1, 1, 1)
	}
}

// mongodb 删除
func BenchmarkDeleteAnswerSheetBySurveyID(b *testing.B) {
	log.ZapInit()
	d := New(mysql.MysqlInit(), mongodb.MongodbInit())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.DeleteAnswerSheetBySurveyID(context.Background(), 1)
	}
}
