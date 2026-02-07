// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

// SessionParticipants is the golang structure for table session_participants.
type SessionParticipants struct {
	Id               int64       `json:"id"               orm:"id"                 description:""`
	SessionId        uuid.UUID   `json:"sessionId"        orm:"session_id"         description:""`
	UserId           uuid.UUID   `json:"userId"           orm:"user_id"            description:""`
	Score            int         `json:"score"            orm:"score"              description:""`
	Rank             int         `json:"rank"             orm:"rank"               description:""`
	CorrectAnswers   int         `json:"correctAnswers"   orm:"correct_answers"    description:""`
	WrongAnswers     int         `json:"wrongAnswers"     orm:"wrong_answers"      description:""`
	TotalTimeSeconds int         `json:"totalTimeSeconds" orm:"total_time_seconds" description:""`
	IsActive         bool        `json:"isActive"         orm:"is_active"          description:""`
	JoinedAt         *gtime.Time `json:"joinedAt"         orm:"joined_at"          description:""`
	LastActivity     *gtime.Time `json:"lastActivity"     orm:"last_activity"      description:""`
}
