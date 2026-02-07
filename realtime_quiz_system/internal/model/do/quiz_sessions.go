// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// QuizSessions is the golang structure of table quiz_sessions for DAO operations like Where/Data.
type QuizSessions struct {
	g.Meta               `orm:"table:quiz_sessions, do:true"`
	Id                   any         //
	QuizId               any         //
	SessionCode          any         //
	HostId               any         //
	Status               any         //
	CurrentQuestionId    any         //
	CurrentQuestionIndex any         //
	MaxParticipants      any         //
	StartedAt            *gtime.Time //
	EndedAt              *gtime.Time //
	CreatedAt            *gtime.Time //
	UpdatedAt            *gtime.Time //
}
