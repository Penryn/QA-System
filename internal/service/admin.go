package service

import (
	"QA-System/internal/dao"
	global "QA-System/internal/global/config"
	"QA-System/internal/models"
	"QA-System/internal/pkg/log"
	"QA-System/internal/pkg/utils"
	"bufio"
	"encoding/json"
	"os"
	"sort"
	"strings"
	"time"
)

func GetAdminByUsername(username string) (*models.User, error) {
	user, err := dao.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if user.Password != "" {
		aesDecryptPassword(user)
	}
	return user, nil
}

func GetAdminByID(id int) (*models.User, error) {
	user, err := dao.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	if user.Password != "" {
		aesDecryptPassword(user)
	}
	return user, nil
}

func IsAdminExist(username string) error {
	_, err := dao.GetUserByUsername(username)
	return err
}

func CreateAdmin(user models.User) error {
	aesEncryptPassword(&user)
	err := dao.CreateUser(user)
	return err
}

func GetUserByName(username string) (*models.User, error) {
	user, err := dao.GetUserByUsername(username)
	return user, err
}

func CreatePermission(id int, surveyID int) error {
	err := dao.CreateManage(id, surveyID)
	return err
}

func DeletePermission(id int, surveyID int) error {
	err := dao.DeleteManage(id, surveyID)
	return err
}

func CheckPermission(id int, surveyID int) error {
	err := dao.CheckManage(id, surveyID)
	return err
}

func CreateSurvey(id int, title string, desc string, img string, questions []dao.Question, status int, time time.Time) error {
	var survey models.Survey
	survey.UserID = id
	survey.Title = title
	survey.Desc = desc
	survey.Img = img
	survey.Status = status
	survey.Deadline = time
	err := dao.CreateSurvey(survey)
	if err != nil {
		return err
	}
	_, err = createQuestionsAndOptions(questions, survey.ID)
	return err
}

func UpdateSurveyStatus(id int, status int) error {
	err := dao.UpdateSurveyStatus(id, status)
	return err
}

func UpdateSurvey(id int, title string, desc string, img string, questions []dao.Question, time time.Time) error {
	//遍历原有问题，删除对应选项
	var oldQuestions []models.Question
	var old_imgs []string
	new_imgs := make([]string, 0)
	//获取原有图片
	oldQuestions, err := dao.GetQuestionsBySurveyID(id)
	if err != nil {
		return err
	}
	old_imgs, err = getOldImgs(id, oldQuestions)
	if err != nil {
		return err
	}
	//删除原有问题和选项
	for _, oldQuestion := range oldQuestions {
		err = dao.DeleteOption(oldQuestion.ID)
		if err != nil {
			return err
		}
	}
	err = dao.DeleteQuestion(id)
	if err != nil {
		return err
	}
	//修改问卷信息
	err = dao.UpdateSurvey(id, title, desc, img, time)
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
	_, err := dao.GetManageByUIDAndSID(uid, sid)
	return err == nil
}

func DeleteSurvey(id int) error {
	var questions []models.Question
	questions, err := dao.GetQuestionsBySurveyID(id)
	if err != nil {
		return err
	}
	var answerSheets []dao.AnswerSheet
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
		err = dao.DeleteOption(question.ID)
		if err != nil {
			return err
		}
	}
	err = dao.DeleteQuestionBySurveyID(id)
	if err != nil {
		return err
	}
	err = dao.DeleteSurvey(id)
	if err != nil {
		return err
	}
	err = dao.DeleteManageBySurveyID(id)
	return err
}

func GetSurveyAnswers(id int, num int, size int) (dao.AnswersResonse, *int64, error) {
	var answerSheets []dao.AnswerSheet
	data := make([]dao.QuestionAnswers, 0)
	time := make([]string, 0)
	var total *int64
	//获取问题
	questions, err := dao.GetQuestionsBySurveyID(id)
	if err != nil {
		return dao.AnswersResonse{}, nil, err
	}
	//初始化data
	for _, question := range questions {
		var q dao.QuestionAnswers
		q.Title = question.Subject
		q.Answers = make([]string, 0)
		data = append(data, q)
	}
	//获取答卷
	answerSheets, total, err = GetAnswerSheetBySurveyID(id, num, size)
	if err != nil {
		return dao.AnswersResonse{}, nil, err
	}
	//填充data
	for _, answerSheet := range answerSheets {
		time = append(time, answerSheet.Time)
		for _, answer := range answerSheet.Answers {
			question, err := dao.GetQuestionByID(answer.QuestionID)
			if err != nil {
				return dao.AnswersResonse{}, nil, err
			}
			for i, q := range data {
				if q.Title == question.Subject {
					data[i].Answers = append(data[i].Answers, answer.Content)
				}
			}
		}
	}
	return dao.AnswersResonse{QuestionAnswers: data, Time: time}, total, nil
}

func GetAllSurveyByUserID(userId int) ([]interface{}, error) {
	surveys, err := dao.GetAllSurveyByUserID(userId)
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

func GetAllSurvey(pageNum, pageSize int, title string) ([]interface{}, *int64, error) {
	surveys, num, error := dao.GetSurveyByTitle(title, pageNum, pageSize)
	if error != nil {
		return nil, nil, error
	}
	response := getSurveyResponse(surveys)
	return response, num, nil
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
	var manages []models.Manage
	manages, err := dao.GetManageByUserID(userId)
	return manages, err
}

func GetAllSurveyAnswers(id int) (dao.AnswersResonse, error) {
	var data []dao.QuestionAnswers
	var answerSheets []dao.AnswerSheet
	var questions []models.Question
	var time []string
	questions, err := dao.GetQuestionsBySurveyID(id)
	if err != nil {
		return dao.AnswersResonse{}, err
	}
	for _, question := range questions {
		var q dao.QuestionAnswers
		q.Title = question.Subject
		data = append(data, q)
	}
	answerSheets, _, err = GetAnswerSheetBySurveyID(id, 0, 0)
	if err != nil {
		return dao.AnswersResonse{}, err
	}
	for _, answerSheet := range answerSheets {
		time = append(time, answerSheet.Time)
		for _, answer := range answerSheet.Answers {
			question, err := dao.GetQuestionByID(answer.QuestionID)
			if err != nil {
				return dao.AnswersResonse{}, err
			}
			for i, q := range data {
				if q.Title == question.Subject {
					data[i].Answers = append(data[i].Answers, answer.Content)
				}
			}
		}
	}
	return dao.AnswersResonse{QuestionAnswers: data, Time: time}, nil
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
	survey, err := dao.GetSurveyByID(id)
	if err != nil {
		return nil, err
	}
	imgs = append(imgs, survey.Img)
	for _, question := range questions {
		imgs = append(imgs, question.Img)
		var options []models.Option
		options, err = dao.GetOptionsByQuestionID(question.ID)
		if err != nil {
			return nil, err
		}
		for _, option := range options {
			imgs = append(imgs, option.Img)
		}
	}
	return imgs, nil
}

func getDelImgs(id int, questions []models.Question, answerSheets []dao.AnswerSheet) ([]string, error) {
	var imgs []string
	survey, err := dao.GetSurveyByID(id)
	if err != nil {
		return nil, err
	}
	imgs = append(imgs, survey.Img)
	for _, question := range questions {
		imgs = append(imgs, question.Img)
		var options []models.Option
		options, err = dao.GetOptionsByQuestionID(question.ID)
		if err != nil {
			return nil, err
		}
		for _, option := range options {
			imgs = append(imgs, option.Img)
		}
	}
	for _, answerSheet := range answerSheets {
		for _, answer := range answerSheet.Answers {
			question, err := dao.GetQuestionByID(answer.QuestionID)
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

func createQuestionsAndOptions(questions []dao.Question, sid int) ([]string, error) {
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
		err := dao.CreateQuestion(q)
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
			err := dao.CreateOption(o)
			if err != nil {
				return nil, err
			}
		}
	}
	return imgs, nil
}

func GetLastLinesFromLogFile(numLines int, logType int) ([]map[string]interface{}, error) {
	levelMap := map[int]string{
		0: "",
		1: "ERROR",
		2: "WARN",
		3: "INFO",
		4: "DEBUG",
	}
	level := levelMap[logType]

	var files []*os.File
	var file *os.File
	var err error

	if logType == 0 {
		// 打开所有相关的日志文件
		files, err = openAllLogFiles()
		if err != nil {
			return nil, err
		}
	} else {
		// 根据 logType 打开特定的日志文件
		file, err = openLogFile(logType)
		if err != nil {
			return nil, err
		}
		if file != nil {
			files = append(files, file)
		}
	}
	defer closeFiles(files)

	if len(files) == 0 {
		return nil, nil
	}

	// 用于存储解析后的日志内容
	var logs []map[string]interface{}

	// 从每个文件中读取内容
	for _, file := range files {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// 解析 JSON 字符串为 map 类型
			var logData map[string]interface{}
			if err := json.Unmarshal(scanner.Bytes(), &logData); err != nil {
				// 如果解析失败，跳过这行日志继续处理下一行
				continue
			}

			// 根据 logType 筛选日志
			if level != "" {
				if logData["L"] == level {
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
	}

	// 如果文件中的行数不足以满足需求，直接返回所有行
	if len(logs) <= numLines {
		return logs, nil
	}

	// 如果文件中的行数超过需求，提取最后几行并返回
	startIndex := len(logs) - numLines
	return logs[startIndex:], nil
}

// 根据 logType 打开单个日志文件
func openLogFile(logType int) (*os.File, error) {
	var filePath string
	switch logType {
	case 1:
		filePath = log.LogDir + "/" + log.LogName + log.ErrorLogSuffix
	case 2:
		filePath = log.LogDir + "/" + log.LogName + log.WarnLogSuffix
	case 3, 4:
		filePath = log.LogDir + "/" + log.LogName + log.LogSuffix
	}
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // 文件不存在，返回 nil
		}
		return nil, err
	}
	return file, nil
}

// 打开所有相关的日志文件
func openAllLogFiles() ([]*os.File, error) {
	filePaths := []string{
		log.LogDir + "/" + log.LogName + log.LogSuffix,
		log.LogDir + "/" + log.LogName + log.ErrorLogSuffix,
		log.LogDir + "/" + log.LogName + log.WarnLogSuffix,
	}

	var openFiles []*os.File
	for _, filePath := range filePaths {
		f, err := os.Open(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			closeFiles(openFiles)
			return nil, err
		}
		openFiles = append(openFiles, f)
	}
	return openFiles, nil
}

// 关闭所有文件
func closeFiles(files []*os.File) {
	for _, file := range files {
		file.Close()
	}
}

func DeleteAnswerSheetBySurveyID(surveyID int) error {
	err := dao.DeleteAnswerSheetBySurveyID(surveyID)
	return err
}

func aesDecryptPassword(user *models.User) {
	user.Password = utils.AesDecrypt(user.Password)
}

func aesEncryptPassword(user *models.User) {
	user.Password = utils.AesEncrypt(user.Password)
}
