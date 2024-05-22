package user

import (
	"QA-System/internal/dao"
	"QA-System/internal/handler/queue"
	"QA-System/internal/pkg/code"
	"QA-System/internal/pkg/queue/asynq"
	"QA-System/internal/pkg/utils"
	"QA-System/internal/service"
	"errors"

	"image/jpeg"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

type SubmitServeyData struct {
	ID            int                         `json:"id" binding:"required"`
	QuestionsList []dao.QuestionsList `json:"questions_list"`
}

func SubmitSurvey(c *gin.Context) {
	var data SubmitServeyData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取参数失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ParamError)
		return
	}
	// 判断问卷问题和答卷问题数目是否一致
	survey, err := service.GetSurveyByID(data.ID)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取问卷失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	questions, err := service.GetQuestionsBySurveyID(survey.ID)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取问题失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	if len(questions) != len(data.QuestionsList) {
		c.Error(&gin.Error{Err: errors.New("问题数量不一致"), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	// 判断填写时间是否在问卷有效期内
	if !survey.Deadline.IsZero() && survey.Deadline.Before(time.Now()) {
		c.Error(&gin.Error{Err: errors.New("填写时间已过"), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.TimeBeyondError)
		return
	}
	// 逐个判断问题答案
	for _, q := range data.QuestionsList {
		question, err := service.GetQuestionByID(q.QuestionID)
		if err != nil {
			c.Error(&gin.Error{Err: errors.New("获取问题失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
			utils.JsonErrorResponse(c, code.ServerError)
			return
		}
		if question.SerialNum != q.SerialNum {
			c.Error(&gin.Error{Err: errors.New("问题序号不一致"), Type: gin.ErrorTypeAny})
			utils.JsonErrorResponse(c, code.ServerError)
			return
		}
		if question.SurveyID != survey.ID {
			c.Error(&gin.Error{Err: errors.New("问题不属于该问卷"), Type: gin.ErrorTypeAny})
			utils.JsonErrorResponse(c, code.ServerError)
			return
		}
		// 判断必填字段是否为空
		if question.Required && q.Answer == "" {
			c.Error(&gin.Error{Err: errors.New("必填字段为空"), Type: gin.ErrorTypeAny})
			utils.JsonErrorResponse(c, code.ServerError)
			return
		}
		// 判断唯一字段是否唯一
		if question.Unique {
			unique, err := service.CheckUnique(data.ID, q.QuestionID, question.SerialNum, q.Answer)
			if err != nil {
				c.Error(&gin.Error{Err: errors.New("唯一字段检查失败"), Type: gin.ErrorTypeAny})
				utils.JsonErrorResponse(c, code.ServerError)
				return
			}
			if !unique {
				c.Error(&gin.Error{Err: errors.New("唯一字段不唯一"), Type: gin.ErrorTypeAny})
				utils.JsonErrorResponse(c, code.UniqueError)
				return
			}

		}
	}
	// 创建并入队任务
    task, err := queue.NewSubmitSurveyTask(data.ID, data.QuestionsList)
    if err != nil {
        c.Error(&gin.Error{Err: errors.New("创建任务失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
        utils.JsonErrorResponse(c, code.ServerError)
        return
    }
	
	_, err = asynq.Client.Enqueue(task)
    if err != nil {
        c.Error(&gin.Error{Err: errors.New("任务入队失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
        utils.JsonErrorResponse(c, code.ServerError)
        return
    }
	utils.JsonSuccessResponse(c, nil)
}

type GetSurveyData struct {
	ID int `form:"id" binding:"required"`
}

type SurveyData struct {
	ID        int                    `json:"id"`
	Time      string                 `json:"time"`
	Desc      string                 `json:"desc"`
	Img       string                 `json:"img"`
	Questions []dao.Question `json:"questions"`
}

// 用户获取问卷
func GetSurvey(c *gin.Context) {
	var data GetSurveyData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取参数失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ParamError)
		return
	}
	// 获取问卷
	survey, err := service.GetSurveyByID(data.ID)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取问卷失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	// 判断填写时间是否在问卷有效期内
	if !survey.Deadline.IsZero() && survey.Deadline.Before(time.Now()) {
		c.Error(&gin.Error{Err: errors.New("填写时间已过"), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.TimeBeyondError)
		return
	}
	// 获取相应的问题
	questions, err := service.GetQuestionsBySurveyID(survey.ID)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取问题失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	// 构建问卷响应
	questionsResponse := make([]map[string]interface{}, 0)
	for _, question := range questions {
		options, err := service.GetOptionsByQuestionID(question.ID)
		if err != nil {
			c.Error(&gin.Error{Err: errors.New("获取选项失败原因: " + err.Error()), Type: gin.ErrorTypeAny})
			utils.JsonErrorResponse(c, code.ServerError)
			return
		}
		optionsResponse := make([]map[string]interface{}, 0)
		for _, option := range options {
			optionResponse := map[string]interface{}{
				"img":        option.Img,
				"content":    option.Content,
				"serial_num": option.SerialNum,
			}
			optionsResponse = append(optionsResponse, optionResponse)
		}
		questionMap := map[string]interface{}{
			"id":            question.ID,
			"serial_num":    question.SerialNum,
			"subject":       question.Subject,
			"describe":      question.Description,
			"required":      question.Required,
			"unique":        question.Unique,
			"other_option":  question.OtherOption,
			"img":           question.Img,
			"question_type": question.QuestionType,
			"reg":           question.Reg,
			"options":       optionsResponse,
		}
		questionsResponse = append(questionsResponse, questionMap)
	}
	response := map[string]interface{}{
		"id":        survey.ID,
		"title":     survey.Title,
		"time":      survey.Deadline,
		"desc":      survey.Desc,
		"img":       survey.Img,
		"questions": questionsResponse,
	}

	utils.JsonSuccessResponse(c, response)
}

// 上传图片
func UploadImg(c *gin.Context) {
	// 保存图片文件
	file, err := c.FormFile("img")
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("获取文件失败"), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	// 检查文件类型是否为图像
	if !isImageFile(file) {
		c.Error(&gin.Error{Err: errors.New("文件类型不是图片"), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.PictureError)
		return
	}
	// 检查文件大小是否超出限制
	if file.Size > 10<<20 { // 10MB，1MB = 1024 * 1024 bytes
		c.Error(&gin.Error{Err: errors.New("文件大小超出限制"), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.PictureSizeError)
		return
	}
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "tempdir")
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("创建临时目录失败"+err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			c.Error(&gin.Error{Err: errors.New("删除临时目录失败"+err.Error()), Type: gin.ErrorTypeAny})
			utils.JsonErrorResponse(c, code.ServerError)
			return
		}
	}() // 在处理完之后删除临时目录及其中的文件
	// 在临时目录中创建临时文件
	tempFile := filepath.Join(tempDir, file.Filename)
	f, err := os.Create(tempFile)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("创建临时文件失败"+err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			c.Error(&gin.Error{Err: errors.New("关闭文件失败"+err.Error()), Type: gin.ErrorTypeAny})
			utils.JsonErrorResponse(c, code.ServerError)
			return
		}
	}()
	// 将上传的文件保存到临时文件中
	src, err := file.Open()
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("打开文件失败"+err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	defer func() {
		if err := src.Close(); err != nil {
			c.Error(&gin.Error{Err: errors.New("关闭文件失败"+err.Error()), Type: gin.ErrorTypeAny})
			utils.JsonErrorResponse(c, code.ServerError)
			return
		}
	}()

	_, err = io.Copy(f, src)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("保存文件失败"+err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}
	// 判断文件的MIME类型是否为图片
	mime, err := mimetype.DetectFile(tempFile)
	if err != nil || !strings.HasPrefix(mime.String(), "image/") {
		c.Error(&gin.Error{Err: errors.New("文件类型不是图片"+err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.PictureError)
		return
	}
	// 保存原始图片
	filename := uuid.New().String() + ".jpg" // 修改扩展名为.jpg
	dst := "./public/static/" + filename
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("保存文件失败"+err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}

	// 转换图像为JPG格式并压缩
	jpgFile := filepath.Join(tempDir, "compressed.jpg")
	err = convertAndCompressImage(dst, jpgFile)
	if err != nil {
		c.Error(&gin.Error{Err: errors.New("转换和压缩图像失败原因:"+err.Error()), Type: gin.ErrorTypeAny})
		utils.JsonErrorResponse(c, code.ServerError)
		return
	}

	//替换原始文件为压缩后的JPG文件
	err = os.Rename(jpgFile, dst)
	if err != nil {
		err = copyFile(jpgFile, dst)
		if err != nil {
			c.Error(&gin.Error{Err: errors.New("替换文件失败原因:"+err.Error()), Type: gin.ErrorTypeAny})
			utils.JsonErrorResponse(c, code.ServerError)
			return
		}
		// Remove the temporary file after copying
		err = os.Remove(jpgFile)
		if err != nil {
			c.Error(&gin.Error{Err: errors.New("删除临时文件失败原因:"+err.Error()), Type: gin.ErrorTypeAny})
			utils.JsonErrorResponse(c, code.ServerError)
			return
		}
	}

	urlHost := service.GetConfigUrl()
	url := urlHost + "/public/static/" + filename

	utils.JsonSuccessResponse(c, url)
}

// 仅支持常见的图像文件类型
func isImageFile(file *multipart.FileHeader) bool {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}
	return allowedTypes[file.Header.Get("Content-Type")]
}

// 用于转换和压缩图像的函数
func convertAndCompressImage(srcPath, dstPath string) error {
	srcImg, err := imaging.Open(srcPath)
	if err != nil {
		return err
	}

	// 调整图像大小（根据需要进行调整）
	resizedImg := resize.Resize(300, 0, srcImg, resize.Lanczos3)

	// 创建新的JPG文件
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 以JPG格式保存调整大小的图像，并设置压缩质量为90
	err = jpeg.Encode(dstFile, resizedImg, &jpeg.Options{Quality: 90})
	if err != nil {
		return err
	}

	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
