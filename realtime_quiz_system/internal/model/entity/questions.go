// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

// Questions is the golang structure for table questions.
type Questions struct {
	Id               int64       `json:"id"               orm:"id"                 description:""`
	QuizId           uuid.UUID   `json:"quizId"           orm:"quiz_id"            description:""`
	QuestionText     string      `json:"questionText"     orm:"question_text"      description:""`
	QuestionType     string      `json:"questionType"     orm:"question_type"      description:""`
	QuestionOrder    int         `json:"questionOrder"    orm:"question_order"     description:""`
	Points           int         `json:"points"           orm:"points"             description:""`
	TimeLimitSeconds int         `json:"timeLimitSeconds" orm:"time_limit_seconds" description:""`
	Explanation      string      `json:"explanation"      orm:"explanation"        description:""`
	CreatedAt        *gtime.Time `json:"createdAt"        orm:"created_at"         description:""`
	UpdatedAt        *gtime.Time `json:"updatedAt"        orm:"updated_at"         description:""`
}
