package queue

import (
	"QA-System/internal/dao"
	"encoding/json"

	"github.com/hibiken/asynq"
)

type SubmitSurveyPayload struct {
    ID            int                  `json:"id"`
    QuestionsList []dao.QuestionsList `json:"questions_list"`
}

const TypeSubmitSurvey = "survey:submit"

func NewSubmitSurveyTask(id int, questionsList []dao.QuestionsList) (*asynq.Task, error) {
    payload, err := json.Marshal(SubmitSurveyPayload{ID: id, QuestionsList: questionsList})
    if err != nil {
        return nil, err
    }
    return asynq.NewTask(TypeSubmitSurvey, payload), nil
}