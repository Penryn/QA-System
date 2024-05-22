package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AnswerSheet struct {
	SurveyID int      `json:"survey_id"` // 问卷ID
	Time     string   `json:"time"`      // 回答时间
	Answers  []Answer `json:"answers"`   // 回答
}

type Answer struct {
	QuestionID int    `json:"question_id"` // 问题ID
	SerialNum  int    `json:"serial_num"`  // 问题序号
	Subject    string `json:"subject"`     // 问题
	Content    string `json:"content"`     // 回答内容
}

type QuestionAnswers struct {
	Title   string   `json:"title"`
	Answers []string `json:"answers"`
}

type AnswersResonse struct {
	QuestionAnswers []QuestionAnswers `json:"question_answers"`
	Time            []string          `json:"time"`
}

// SaveAnswerSheet 将答卷保存到 MongoDB 集合中
func (d *Dao) SaveAnswerSheet(ctx context.Context, answerSheet AnswerSheet) error {
	_, err := d.mongo.InsertOne(ctx, answerSheet)
	return err
}

// GetAnswerSheetBySurveyID 根据问卷ID分页获取答卷
func (d *Dao) GetAnswerSheetBySurveyID(ctx context.Context, surveyID int, pageNum int, pageSize int) ([]AnswerSheet, *int64, error) {
	var answerSheets []AnswerSheet
	filter := bson.M{"surveyid": surveyID}

	// 设置总记录数查询过滤条件和选项
	countFilter := bson.M{"surveyid": surveyID}
	countOpts := options.Count()

	// 执行总记录数查询
	total, err := d.mongo.CountDocuments(ctx, countFilter, countOpts)
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
	cur, err := d.mongo.Find(ctx, filter, opts)
	if err != nil {
		return nil, nil, err
	}
	defer cur.Close(ctx)

	// 迭代查询结果
	for cur.Next(ctx) {
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

// DeleteAnswerSheetBySurveyID 根据问卷ID删除答卷
func (d *Dao) DeleteAnswerSheetBySurveyID(ctx context.Context, surveyID int) error {
	filter := bson.M{"surveyid": surveyID}
	// 删除所有满足条件的文档
	_, err := d.mongo.DeleteMany(ctx, filter)
	return err
}
