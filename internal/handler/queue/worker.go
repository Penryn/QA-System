package queue

import (
	"QA-System/internal/service"
	"context"
	"encoding/json"
	"errors"

	"github.com/hibiken/asynq"
)

func HandleSubmitSurveyTask(ctx context.Context, t *asynq.Task) error {
	var p SubmitSurveyPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	// 提交问卷
	err := service.SubmitSurvey(p.ID, p.QuestionsList)
	if err != nil {
		return errors.New("提交问卷失败原因: " + err.Error())
	}

	return nil
}