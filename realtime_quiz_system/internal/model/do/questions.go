// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Questions is the golang structure of table questions for DAO operations like Where/Data.
type Questions struct {
	g.Meta           `orm:"table:questions, do:true"`
	Id               any         //
	QuizId           any         //
	QuestionText     any         //
	QuestionType     any         //
	QuestionOrder    any         //
	Points           any         //
	TimeLimitSeconds any         //
	Explanation      any         //
	CreatedAt        *gtime.Time //
	UpdatedAt        *gtime.Time //
}
