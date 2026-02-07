// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SessionParticipants is the golang structure of table session_participants for DAO operations like Where/Data.
type SessionParticipants struct {
	g.Meta           `orm:"table:session_participants, do:true"`
	Id               any         //
	SessionId        any         //
	UserId           any         //
	Score            any         //
	Rank             any         //
	CorrectAnswers   any         //
	WrongAnswers     any         //
	TotalTimeSeconds any         //
	IsActive         any         //
	JoinedAt         *gtime.Time //
	LastActivity     *gtime.Time //
}
