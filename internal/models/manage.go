package models

type Manage struct {
	ID       int `json:"id"`
	UserID   int `json:"user_id"`
	SurveyID int `json:"survey_id"`
}
