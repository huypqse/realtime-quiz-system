// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

// QuizSessions is the golang structure for table quiz_sessions.
type QuizSessions struct {
	Id                   uuid.UUID   `json:"id"                   orm:"id"                     description:""`
	QuizId               uuid.UUID   `json:"quizId"               orm:"quiz_id"                description:""`
	SessionCode          string      `json:"sessionCode"          orm:"session_code"           description:""`
	HostId               uuid.UUID   `json:"hostId"               orm:"host_id"                description:""`
	Status               string      `json:"status"               orm:"status"                 description:""`
	CurrentQuestionId    int64       `json:"currentQuestionId"    orm:"current_question_id"    description:""`
	CurrentQuestionIndex int         `json:"currentQuestionIndex" orm:"current_question_index" description:""`
	MaxParticipants      int         `json:"maxParticipants"      orm:"max_participants"       description:""`
	StartedAt            *gtime.Time `json:"startedAt"            orm:"started_at"             description:""`
	EndedAt              *gtime.Time `json:"endedAt"              orm:"ended_at"               description:""`
	CreatedAt            *gtime.Time `json:"createdAt"            orm:"created_at"             description:""`
	UpdatedAt            *gtime.Time `json:"updatedAt"            orm:"updated_at"             description:""`
}
