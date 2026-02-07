// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// LeaderboardSnapshotsDao is the data access object for the table leaderboard_snapshots.
type LeaderboardSnapshotsDao struct {
	table    string                      // table is the underlying table name of the DAO.
	group    string                      // group is the database configuration group name of the current DAO.
	columns  LeaderboardSnapshotsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler          // handlers for customized model modification.
}

// LeaderboardSnapshotsColumns defines and stores column names for the table leaderboard_snapshots.
type LeaderboardSnapshotsColumns struct {
	Id           string //
	SessionId    string //
	UserId       string //
	Score        string //
	Rank         string //
	SnapshotType string //
	QuestionId   string //
	CreatedAt    string //
}

// leaderboardSnapshotsColumns holds the columns for the table leaderboard_snapshots.
var leaderboardSnapshotsColumns = LeaderboardSnapshotsColumns{
	Id:           "id",
	SessionId:    "session_id",
	UserId:       "user_id",
	Score:        "score",
	Rank:         "rank",
	SnapshotType: "snapshot_type",
	QuestionId:   "question_id",
	CreatedAt:    "created_at",
}

// NewLeaderboardSnapshotsDao creates and returns a new DAO object for table data access.
func NewLeaderboardSnapshotsDao(handlers ...gdb.ModelHandler) *LeaderboardSnapshotsDao {
	return &LeaderboardSnapshotsDao{
		group:    "default",
		table:    "leaderboard_snapshots",
		columns:  leaderboardSnapshotsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *LeaderboardSnapshotsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *LeaderboardSnapshotsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *LeaderboardSnapshotsDao) Columns() LeaderboardSnapshotsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *LeaderboardSnapshotsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *LeaderboardSnapshotsDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *LeaderboardSnapshotsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
