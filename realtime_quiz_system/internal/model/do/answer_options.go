// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AnswerOptions is the golang structure of table answer_options for DAO operations like Where/Data.
type AnswerOptions struct {
	g.Meta      `orm:"table:answer_options, do:true"`
	Id          any         //
	QuestionId  any         //
	OptionText  any         //
	IsCorrect   any         //
	OptionOrder any         //
	CreatedAt   *gtime.Time //
}
