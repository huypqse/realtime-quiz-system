// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Quizzes is the golang structure of table quizzes for DAO operations like Where/Data.
type Quizzes struct {
	g.Meta           `orm:"table:quizzes, do:true"`
	Id               any         //
	Title            any         //
	Description      any         //
	CreatedBy        any         //
	Status           any         //
	StartTime        *gtime.Time //
	EndTime          *gtime.Time //
	TimeLimitSeconds any         //
	TotalQuestions   any         //
	CreatedAt        *gtime.Time //
	UpdatedAt        *gtime.Time //
}
