// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserAnswers is the golang structure of table user_answers for DAO operations like Where/Data.
type UserAnswers struct {
	g.Meta           `orm:"table:user_answers, do:true"`
	Id               any         //
	SessionId        any         //
	UserId           any         //
	QuestionId       any         //
	SelectedOptionId any         //
	AnswerText       any         //
	IsCorrect        any         //
	PointsEarned     any         //
	TimeTakenSeconds any         //
	AnsweredAt       *gtime.Time //
}
