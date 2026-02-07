// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

// UserAnswers is the golang structure for table user_answers.
type UserAnswers struct {
	Id               int64       `json:"id"               orm:"id"                 description:""`
	SessionId        uuid.UUID   `json:"sessionId"        orm:"session_id"         description:""`
	UserId           uuid.UUID   `json:"userId"           orm:"user_id"            description:""`
	QuestionId       int64       `json:"questionId"       orm:"question_id"        description:""`
	SelectedOptionId int64       `json:"selectedOptionId" orm:"selected_option_id" description:""`
	AnswerText       string      `json:"answerText"       orm:"answer_text"        description:""`
	IsCorrect        bool        `json:"isCorrect"        orm:"is_correct"         description:""`
	PointsEarned     int         `json:"pointsEarned"     orm:"points_earned"      description:""`
	TimeTakenSeconds int         `json:"timeTakenSeconds" orm:"time_taken_seconds" description:""`
	AnsweredAt       *gtime.Time `json:"answeredAt"       orm:"answered_at"        description:""`
}
