// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// QuizSessionsDao is the data access object for the table quiz_sessions.
type QuizSessionsDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  QuizSessionsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// QuizSessionsColumns defines and stores column names for the table quiz_sessions.
type QuizSessionsColumns struct {
	Id                   string //
	QuizId               string //
	SessionCode          string //
	HostId               string //
	Status               string //
	CurrentQuestionId    string //
	CurrentQuestionIndex string //
	MaxParticipants      string //
	StartedAt            string //
	EndedAt              string //
	CreatedAt            string //
	UpdatedAt            string //
}

// quizSessionsColumns holds the columns for the table quiz_sessions.
var quizSessionsColumns = QuizSessionsColumns{
	Id:                   "id",
	QuizId:               "quiz_id",
	SessionCode:          "session_code",
	HostId:               "host_id",
	Status:               "status",
	CurrentQuestionId:    "current_question_id",
	CurrentQuestionIndex: "current_question_index",
	MaxParticipants:      "max_participants",
	StartedAt:            "started_at",
	EndedAt:              "ended_at",
	CreatedAt:            "created_at",
	UpdatedAt:            "updated_at",
}

// NewQuizSessionsDao creates and returns a new DAO object for table data access.
func NewQuizSessionsDao(handlers ...gdb.ModelHandler) *QuizSessionsDao {
	return &QuizSessionsDao{
		group:    "default",
		table:    "quiz_sessions",
		columns:  quizSessionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *QuizSessionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *QuizSessionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *QuizSessionsDao) Columns() QuizSessionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *QuizSessionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *QuizSessionsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *QuizSessionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
