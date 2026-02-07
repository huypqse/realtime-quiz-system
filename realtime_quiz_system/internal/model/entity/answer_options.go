// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AnswerOptions is the golang structure for table answer_options.
type AnswerOptions struct {
	Id          int64       `json:"id"          orm:"id"           description:""`
	QuestionId  int64       `json:"questionId"  orm:"question_id"  description:""`
	OptionText  string      `json:"optionText"  orm:"option_text"  description:""`
	IsCorrect   bool        `json:"isCorrect"   orm:"is_correct"   description:""`
	OptionOrder int         `json:"optionOrder" orm:"option_order" description:""`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:""`
}
