// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

// LeaderboardSnapshots is the golang structure for table leaderboard_snapshots.
type LeaderboardSnapshots struct {
	Id           int64       `json:"id"           orm:"id"            description:""`
	SessionId    uuid.UUID   `json:"sessionId"    orm:"session_id"    description:""`
	UserId       uuid.UUID   `json:"userId"       orm:"user_id"       description:""`
	Score        int         `json:"score"        orm:"score"         description:""`
	Rank         int         `json:"rank"         orm:"rank"          description:""`
	SnapshotType string      `json:"snapshotType" orm:"snapshot_type" description:""`
	QuestionId   int64       `json:"questionId"   orm:"question_id"   description:""`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:""`
}
