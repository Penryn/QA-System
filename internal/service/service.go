package service

import (
	global "QA-System/internal/global/config"
	"QA-System/internal/models"
	mongodb "QA-System/internal/pkg/database/mongodb"
	mysql "QA-System/internal/pkg/database/mysql"
	"QA-System/internal/pkg/log"
	"QA-System/internal/pkg/utils"
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type Option struct {
	SerialNum int    `json:"serial_num"` //选项序号
	Content   string `json:"content"`    //选项内容
	Img       string `json:"img"`        //图片
}

type Question struct {
	ID           int      `json:"id"`
	SerialNum    int      `json:"serial_num"`    //题目序号
	Subject      string   `json:"subject"`       //问题
	Description  string   `json:"description"`   //问题描述
	Img          string   `json:"img"`           //图片
	Required     bool     `json:"required"`      //是否必填
	Unique       bool     `json:"unique"`        //是否唯一
	OtherOption  bool     `json:"other_option"`  //是否有其他选项
	QuestionType int      `json:"question_type"` //问题类型 1单选2多选3填空4简答5图片
	Reg          string   `json:"reg"`           //正则表达式
	Options      []Option `json:"options"`       //选项
}

type QuestionsList struct {
	QuestionID int    `json:"question_id" binding:"required"`
	SerialNum  int    `json:"serial_num"`
	Answer     string `json:"answer"`
}

type Answer struct {
	QuestionID int    `json:"question_id"` //问题ID
	SerialNum  int    `json:"serial_num"`  //问题序号
	Subject    string `json:"subject"`     //问题
	Content    string `json:"content"`     //回答内容
}

type AnswerSheet struct {
	SurveyID int      `json:"survey_id"` //问卷ID
	Time     string   `json:"time"`      //回答时间
	Answers  []Answer `json:"answers"`   //回答
}

func GetAdminByUsername(username string) (*models.User, error) {
	var user models.User
	result := mysql.DB.Model(&models.User{}).Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if user.Password != "" {
		aesDecryptPassword(&user)
	}
	return &user, result.Error
}

func GetAdminByID(id int) (*models.User, error) {
	user := models.User{}
	result := mysql.DB.Model(&models.User{}).Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	aesDecryptPassword(&user)
	return &user, nil
}

func IsAdminExist(username string) error {
	var user models.User
	result := mysql.DB.Model(models.User{}).Where("username = ?", username).First(&user)
	return result.Error
}

func CreateAdmin(user models.User) error {
	aesEncryptPassword(&user)
	result := mysql.DB.Model(models.User{}).Create(&user)
	return result.Error
}

func aesDecryptPassword(user *models.User) {
	user.Password = utils.AesDecrypt(user.Password)
}

func aesEncryptPassword(user *models.User) {
	user.Password = utils.AesEncrypt(user.Password)
}

func GetLastLinesFromLogFile(filePath string, numLines int, logType int) ([]map[string]interface{}, error) {
	// 打开日志文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	levelMap := map[int]string{
		0: "",
		1: "error",
		2: "warn",
		3: "info",
		4: "debug",
	}
	level := levelMap[logType]

	// 用于存储解析后的日志内容
	var logs []map[string]interface{}

	// 逐行读取文件内容
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// 解析JSON字符串为map类型
		var logData map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &logData); err != nil {
			// 如果解析失败，跳过这行日志继续处理下一行
			continue
		}

		// 根据logType筛选日志
		if level != "" {
			if logData["level"] == level {
				logs = append(logs, logData)
			}
		} else {
			logs = append(logs, logData)
		}

	}

	// 检查是否发生了读取错误
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// 如果文件中的行数不足以满足需求，直接返回所有行
	if len(logs) <= numLines {
		return logs, nil
	}

	// 如果文件中的行数超过需求，提取最后几行并返回
	startIndex := len(logs) - numLines
	return logs[startIndex:], nil
}

func GetUserByName(username string) (models.User, error) {
	var user models.User
	err := mysql.DB.Where("username = ?", username).First(&user).Error
	return user, err
}

func CreatePermission(id int, surveyID int) error {
	err := mysql.DB.Create(&models.Manage{UserID: id, SurveyID: surveyID}).Error
	return err
}

func DeletePermission(id int, surveyID int) error {
	err := mysql.DB.Where("user_id = ? AND survey_id = ?", id, surveyID).Delete(&models.Manage{}).Error
	return err
}

func CheckPermission(id int, surveyID int) error {
	var manage models.Manage
	err := mysql.DB.Where("user_id = ? AND survey_id = ?", id, surveyID).First(&manage).Error
	return err
}

func CreateSurvey(id int, title string, desc string, img string, questions []Question, status int, time time.Time) error {
	var survey models.Survey
	survey.UserID = id
	survey.Title = title
	survey.Desc = desc
	survey.Img = img
	survey.Status = status
	survey.Deadline = time
	err := mysql.DB.Create(&survey).Error
	if err != nil {
		return err
	}
	_, err = createQuestionsAndOptions(questions, survey.ID)
	return err
}

func UpdateSurveyStatus(id int, status int) error {
	var survey models.Survey
	err := mysql.DB.Model(&survey).Where("id = ?", id).Update("status", status).Error
	return err
}

func UpdateSurvey(id int, title string, desc string, img string, questions []Question, time time.Time) error {
	//遍历原有问题，删除对应选项
	var survey models.Survey
	var oldQuestions []models.Question
	var old_imgs []string
	new_imgs := make([]string, 0)
	//获取原有图片
	err := mysql.DB.Where("survey_id = ?", id).Find(&oldQuestions).Error
	if err != nil {
		return err
	}
	old_imgs, err = getOldImgs(id, oldQuestions)
	if err != nil {
		return err
	}
	//删除原有问题和选项
	for _, oldQuestion := range oldQuestions {
		err = mysql.DB.Where("question_id = ?", oldQuestion.ID).Delete(&models.Option{}).Error
		if err != nil {
			return err
		}
	}
	err = mysql.DB.Where("survey_id = ?", id).Delete(&models.Question{}).Error
	if err != nil {
		return err
	}
	//修改问卷信息
	err = mysql.DB.Model(&survey).Where("id = ?", id).Updates(map[string]interface{}{"title": title, "desc": desc, "img": img, "deadline": time}).Error
	if err != nil {
		return err
	}
	new_imgs = append(new_imgs, img)
	//重新添加问题和选项
	imgs, err := createQuestionsAndOptions(questions, id)
	if err != nil {
		return err
	}
	new_imgs = append(new_imgs, imgs...)
	urlHost := global.Config.GetString("url.host")
	//删除无用图片
	for _, old_img := range old_imgs {
		if !contains(new_imgs, old_img) {
			_ = os.Remove("./static/" + strings.TrimPrefix(old_img, urlHost+"/static/"))
		}
	}
	return nil
}

func UserInManage(uid int, sid int) bool {
	var survey models.Manage
	err := mysql.DB.Where("user_id = ? and survey_id = ?", uid, sid).First(&survey).Error
	return err == nil
}

func DeleteSurvey(id int) error {
	var survey models.Survey
	var questions []models.Question
	err := mysql.DB.Where("survey_id = ?", id).Find(&questions).Error
	if err != nil {
		return err
	}
	var answerSheets []AnswerSheet
	answerSheets, _, err = GetAnswerSheetBySurveyID(id, 0, 0)
	if err != nil {
		return err
	}
	//删除图片
	imgs, err := getDelImgs(id, questions, answerSheets)
	if err != nil {
		return err
	}
	urlHost := global.Config.GetString("url.host")
	for _, img := range imgs {
		_ = os.Remove("./static/" + strings.TrimPrefix(img, urlHost+"/static/"))
	}
	//删除答卷
	err = DeleteAnswerSheetBySurveyID(id)
	if err != nil {
		return err
	}
	//删除问题、选项、问卷、管理
	for _, question := range questions {
		err = mysql.DB.Where("question_id = ?", question.ID).Delete(&models.Option{}).Error
		if err != nil {
			return err
		}
	}
	err = mysql.DB.Where("survey_id = ?", id).Delete(&models.Question{}).Error
	if err != nil {
		return err
	}
	err = mysql.DB.Where("id = ?", id).Delete(&survey).Error
	if err != nil {
		return err
	}
	err = mysql.DB.Where("survey_id = ?", id).Delete(&models.Manage{}).Error
	return err
}

type QuestionAnswers struct {
	Title   string   `json:"title"`
	Answers []string `json:"answers"`
}

type AnswersResonse struct {
	QuestionAnswers []QuestionAnswers `json:"question_answers"`
	Time            []string          `json:"time"`
}

func GetSurveyAnswers(id int, num int, size int) (AnswersResonse, *int64, error) {
	var answerSheets []AnswerSheet
	var questions []models.Question
	data := make([]QuestionAnswers, 0)
	time := make([]string, 0)
	var total *int64
	//获取问题
	err := mysql.DB.Where("survey_id = ?", id).Find(&questions).Error
	if err != nil {
		return AnswersResonse{}, nil, err
	}
	//初始化data
	for _, question := range questions {
		var q QuestionAnswers
		q.Title = question.Subject
		q.Answers = make([]string, 0)
		data = append(data, q)
	}
	//获取答卷
	answerSheets, total, err = GetAnswerSheetBySurveyID(id, num, size)
	if err != nil {
		return AnswersResonse{}, nil, err
	}
	//填充data
	for _, answerSheet := range answerSheets {
		time = append(time, answerSheet.Time)
		for _, answer := range answerSheet.Answers {
			var question models.Question
			err = mysql.DB.Where("id = ?", answer.QuestionID).First(&question).Error
			if err != nil {
				return AnswersResonse{}, nil, err
			}
			for i, q := range data {
				if q.Title == question.Subject {
					data[i].Answers = append(data[i].Answers, answer.Content)
				}
			}
		}
	}
	return AnswersResonse{QuestionAnswers: data, Time: time}, total, nil
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func getOldImgs(id int, questions []models.Question) ([]string, error) {
	var imgs []string
	var survey models.Survey
	err := mysql.DB.Where("id = ?", id).First(&survey).Error
	if err != nil {
		return nil, err
	}
	imgs = append(imgs, survey.Img)
	for _, question := range questions {
		imgs = append(imgs, question.Img)
		var options []models.Option
		err = mysql.DB.Where("question_id = ?", question.ID).Find(&options).Error
		if err != nil {
			return nil, err
		}
		for _, option := range options {
			imgs = append(imgs, option.Img)
		}
	}
	return imgs, nil
}

func getDelImgs(id int, questions []models.Question, answerSheets []AnswerSheet) ([]string, error) {
	var imgs []string
	var survey models.Survey
	err := mysql.DB.Where("id = ?", id).First(&survey).Error
	if err != nil {
		return nil, err
	}
	imgs = append(imgs, survey.Img)
	for _, question := range questions {
		imgs = append(imgs, question.Img)
		var options []models.Option
		err = mysql.DB.Where("question_id = ?", question.ID).Find(&options).Error
		if err != nil {
			return nil, err
		}
		for _, option := range options {
			imgs = append(imgs, option.Img)
		}
	}
	for _, answerSheet := range answerSheets {
		for _, answer := range answerSheet.Answers {
			var question models.Question
			err = mysql.DB.Where("id = ?", answer.QuestionID).First(&question).Error
			if err != nil {
				return nil, err
			}
			if question.QuestionType == 5 {
				imgs = append(imgs, answer.Content)
			}
		}

	}
	return imgs, nil
}

func createQuestionsAndOptions(questions []Question, sid int) ([]string, error) {
	var imgs []string
	for _, question := range questions {
		var q models.Question
		q.SerialNum = question.SerialNum
		q.SurveyID = sid
		q.Subject = question.Subject
		q.Description = question.Description
		q.Img = question.Img
		q.Required = question.Required
		q.Unique = question.Unique
		q.OtherOption = question.OtherOption
		q.QuestionType = question.QuestionType
		imgs = append(imgs, question.Img)
		err := mysql.DB.Create(&q).Error
		if err != nil {
			return nil, err
		}
		for _, option := range question.Options {
			var o models.Option
			o.Content = option.Content
			o.QuestionID = q.ID
			o.SerialNum = option.SerialNum
			o.Img = option.Img
			imgs = append(imgs, option.Img)
			err := mysql.DB.Create(&o).Error
			if err != nil {
				return nil, err
			}
		}
	}
	return imgs, nil
}

func GetAllSurveyByUserID(userId int) ([]interface{}, error) {
	var surveys []models.Survey
	err := mysql.DB.Model(models.Survey{}).Where("user_id = ?", userId).
		Order("CASE WHEN status = 2 THEN 0 ELSE 1 END, id DESC").Find(&surveys).Error
	response := getSurveyResponse(surveys)
	return response, err
}

func ProcessResponse(response []interface{}, pageNum, pageSize int, title string) ([]interface{}, *int64) {
	if title != "" {
		filteredResponse := make([]interface{}, 0)
		for _, item := range response {
			itemMap := item.(map[string]interface{})
			if strings.Contains(strings.ToLower(itemMap["title"].(string)), strings.ToLower(title)) {
				filteredResponse = append(filteredResponse, item)
			}
		}
		response = filteredResponse
	}
	num := int64(len(response))
	sort.Slice(response, func(i, j int) bool {
		return response[i].(map[string]interface{})["id"].(int) > response[j].(map[string]interface{})["id"].(int)
	})
	var sortedResponse []interface{}
	var status2Response, status1Response []interface{}
	for _, item := range response {
		itemMap := item.(map[string]interface{})
		if itemMap["status"].(int) == 2 {
			status2Response = append(status2Response, item)
		} else {
			status1Response = append(status1Response, item)
		}
	}
	sortedResponse = append(status2Response, status1Response...)

	startIdx := (pageNum - 1) * pageSize
	endIdx := startIdx + pageSize
	if endIdx > len(sortedResponse) {
		endIdx = len(sortedResponse)
	}
	pagedResponse := sortedResponse[startIdx:endIdx]

	return pagedResponse, &num
}

func GetAllSurvey(pageNum, pageSize int, title string) ([]interface{}, *int64) {
	var surveys []models.Survey
	var num int64
	query := mysql.DB.Model(models.Survey{}).
		Order("CASE WHEN status = 2 THEN 0 ELSE 1 END, id DESC")
	if title != "" {
		title := "%" + title + "%"
		query = query.Where("title LIKE ?", title)
		query.Find(&surveys)
		num = int64(len(surveys))
	} else {
		query.Find(&surveys)
		num = int64(len(surveys))
	}
	response := getSurveyResponse(surveys)

	startIdx := (pageNum - 1) * pageSize
	endIdx := startIdx + pageSize
	if endIdx > len(response) {
		endIdx = len(response)
	}
	pagedResponse := response[startIdx:endIdx]

	return pagedResponse, &num
}

func getSurveyResponse(surveys []models.Survey) []interface{} {
	response := make([]interface{}, 0)
	for _, survey := range surveys {
		surveyResponse := map[string]interface{}{
			"id":     survey.ID,
			"title":  survey.Title,
			"status": survey.Status,
			"num":    survey.Num,
		}
		response = append(response, surveyResponse)
	}
	return response
}

func GetManageredSurveyByUserID(userId int) ([]models.Manage, error) {
	var surveys []models.Manage
	err := mysql.DB.Model(models.Manage{}).Where("user_id = ?", userId).Order("id DESC").Find(&surveys).Error
	return surveys, err
}

func GetAllSurveyAnswers(id int) (AnswersResonse, error) {
	var data []QuestionAnswers
	var answerSheets []AnswerSheet
	var questions []models.Question
	var time []string
	err := mysql.DB.Where("survey_id = ?", id).Find(&questions).Error
	if err != nil {
		return AnswersResonse{}, err
	}
	for _, question := range questions {
		var q QuestionAnswers
		q.Title = question.Subject
		data = append(data, q)
	}
	answerSheets, _, err = GetAnswerSheetBySurveyID(id, 0, 0)
	if err != nil {
		return AnswersResonse{}, err
	}
	for _, answerSheet := range answerSheets {
		time = append(time, answerSheet.Time)
		for _, answer := range answerSheet.Answers {
			var question models.Question
			err = mysql.DB.Where("id = ?", answer.QuestionID).First(&question).Error
			if err != nil {
				return AnswersResonse{}, err
			}
			for i, q := range data {
				if q.Title == question.Subject {
					data[i].Answers = append(data[i].Answers, answer.Content)
				}
			}
		}
	}
	return AnswersResonse{QuestionAnswers: data, Time: time}, nil
}

func GetSurveyByID(id int) (models.Survey, error) {
	var survey models.Survey
	err := mysql.DB.Where("id = ?", id).First(&survey).Error
	return survey, err
}

func GetQuestionsBySurveyID(id int) ([]models.Question, error) {
	var questions []models.Question
	err := mysql.DB.Where("survey_id = ?", id).Find(&questions).Error
	return questions, err
}

func GetOptionsByQuestionID(questionId int) ([]models.Option, error) {
	var options []models.Option
	err := mysql.DB.Where("question_id = ?", questionId).Find(&options).Error
	return options, err
}

func GetQuestionByID(id int) (models.Question, error) {
	var question models.Question
	err := mysql.DB.Where("id = ?", id).First(&question).Error
	return question, err
}

func CheckUnique(sid int, qid int, serial_num int, content string) (bool, error) {
	var answerSheets []AnswerSheet
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

func SubmitSurvey(sid int, data []QuestionsList) error {
	var answerSheet AnswerSheet
	answerSheet.SurveyID = sid
	answerSheet.Time = time.Now().Format("2006-01-02 15:04:05")
	for _, q := range data {
		var answer Answer
		answer.QuestionID = q.QuestionID
		answer.SerialNum = q.SerialNum
		answer.Content = q.Answer
		answerSheet.Answers = append(answerSheet.Answers, answer)
	}
	err := SaveAnswerSheet(answerSheet)
	if err != nil {
		return err
	}
	err = mysql.DB.Model(&models.Survey{}).Where("id = ?", sid).Update("num", gorm.Expr("num + ?", 1)).Error
	if err != nil {
		return err
	}
	return nil
}

func SaveAnswerSheet(answerSheet AnswerSheet) error {
	_, err := mongodb.MDB.InsertOne(context.Background(), answerSheet)
	if err != nil {
		log.Logger.Error(err.Error())
	}
	return nil
}

func GetAnswerSheetBySurveyID(surveyID int, pageNum int, pageSize int) ([]AnswerSheet, *int64, error) {
	var answerSheets []AnswerSheet
	filter := bson.M{"surveyid": surveyID}

	// 设置总记录数查询过滤条件
	countFilter := bson.M{"surveyid": surveyID}

	// 设置总记录数查询选项
	countOpts := options.Count()

	// 执行总记录数查询
	total, err := mongodb.MDB.CountDocuments(context.Background(), countFilter, countOpts)
	if err != nil {
		return nil, nil, err
	}

	// 设置分页查询选项
	opts := options.Find()
	if pageNum != 0 && pageSize != 0 {
		opts.SetSkip(int64((pageNum - 1) * pageSize)) // 计算要跳过的文档数
		opts.SetLimit(int64(pageSize))                // 设置返回的文档数
	}
	// 执行分页查询
	cur, err := mongodb.MDB.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, nil, err
	}
	defer cur.Close(context.Background())

	// 迭代查询结果
	for cur.Next(context.Background()) {
		var answerSheet AnswerSheet
		if err := cur.Decode(&answerSheet); err != nil {
			return nil, nil, err
		}
		answerSheets = append(answerSheets, answerSheet)
	}
	if err := cur.Err(); err != nil {
		return nil, nil, err
	}

	// 返回分页数据和总记录数
	return answerSheets, &total, nil
}

func DeleteAnswerSheetBySurveyID(surveyID int) error {
	filter := bson.M{"surveyid": surveyID}
	// 删除所有满足条件的文档
	_, err := mongodb.MDB.DeleteMany(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

func SetUserSession(c *gin.Context, user *models.User) error {
	webSession := sessions.Default(c)
	webSession.Options(sessions.Options{
		MaxAge:   3600 * 24 * 7,
		Path:     "/",
		HttpOnly: true,
	})
	webSession.Set("id", user.ID)
	return webSession.Save()
}

func GetUserSession(c *gin.Context) (*models.User, error) {
	webSession := sessions.Default(c)
	id := webSession.Get("id")
	if id == nil {
		return nil, errors.New("")
	}
	user, _ := GetAdminByID(id.(int))
	if user == nil {
		ClearUserSession(c)
		return nil, errors.New("")
	}
	return user, nil
}

func UpdateUserSession(c *gin.Context) (*models.User, error) {
	user, err := GetUserSession(c)
	if err != nil {
		return nil, err
	}
	err = SetUserSession(c, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CheckUserSession(c *gin.Context) bool {
	webSession := sessions.Default(c)
	id := webSession.Get("id")
	return id != nil
}

func ClearUserSession(c *gin.Context) {
	webSession := sessions.Default(c)
	webSession.Delete("id")
	webSession.Save()
}
