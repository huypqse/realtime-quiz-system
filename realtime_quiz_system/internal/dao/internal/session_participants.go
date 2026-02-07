// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SessionParticipantsDao is the data access object for the table session_participants.
type SessionParticipantsDao struct {
	table    string                     // table is the underlying table name of the DAO.
	group    string                     // group is the database configuration group name of the current DAO.
	columns  SessionParticipantsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler         // handlers for customized model modification.
}

// SessionParticipantsColumns defines and stores column names for the table session_participants.
type SessionParticipantsColumns struct {
	Id               string //
	SessionId        string //
	UserId           string //
	Score            string //
	Rank             string //
	CorrectAnswers   string //
	WrongAnswers     string //
	TotalTimeSeconds string //
	IsActive         string //
	JoinedAt         string //
	LastActivity     string //
}

// sessionParticipantsColumns holds the columns for the table session_participants.
var sessionParticipantsColumns = SessionParticipantsColumns{
	Id:               "id",
	SessionId:        "session_id",
	UserId:           "user_id",
	Score:            "score",
	Rank:             "rank",
	CorrectAnswers:   "correct_answers",
	WrongAnswers:     "wrong_answers",
	TotalTimeSeconds: "total_time_seconds",
	IsActive:         "is_active",
	JoinedAt:         "joined_at",
	LastActivity:     "last_activity",
}

// NewSessionParticipantsDao creates and returns a new DAO object for table data access.
func NewSessionParticipantsDao(handlers ...gdb.ModelHandler) *SessionParticipantsDao {
	return &SessionParticipantsDao{
		group:    "default",
		table:    "session_participants",
		columns:  sessionParticipantsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SessionParticipantsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SessionParticipantsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SessionParticipantsDao) Columns() SessionParticipantsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SessionParticipantsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SessionParticipantsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SessionParticipantsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
