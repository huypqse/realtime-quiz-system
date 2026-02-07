// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// LeaderboardSnapshots is the golang structure of table leaderboard_snapshots for DAO operations like Where/Data.
type LeaderboardSnapshots struct {
	g.Meta       `orm:"table:leaderboard_snapshots, do:true"`
	Id           any         //
	SessionId    any         //
	UserId       any         //
	Score        any         //
	Rank         any         //
	SnapshotType any         //
	QuestionId   any         //
	CreatedAt    *gtime.Time //
}
