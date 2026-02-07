// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// QuizzesDao is the data access object for the table quizzes.
type QuizzesDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  QuizzesColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// QuizzesColumns defines and stores column names for the table quizzes.
type QuizzesColumns struct {
	Id               string //
	Title            string //
	Description      string //
	CreatedBy        string //
	Status           string //
	StartTime        string //
	EndTime          string //
	TimeLimitSeconds string //
	TotalQuestions   string //
	CreatedAt        string //
	UpdatedAt        string //
}

// quizzesColumns holds the columns for the table quizzes.
var quizzesColumns = QuizzesColumns{
	Id:               "id",
	Title:            "title",
	Description:      "description",
	CreatedBy:        "created_by",
	Status:           "status",
	StartTime:        "start_time",
	EndTime:          "end_time",
	TimeLimitSeconds: "time_limit_seconds",
	TotalQuestions:   "total_questions",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
}

// NewQuizzesDao creates and returns a new DAO object for table data access.
func NewQuizzesDao(handlers ...gdb.ModelHandler) *QuizzesDao {
	return &QuizzesDao{
		group:    "default",
		table:    "quizzes",
		columns:  quizzesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *QuizzesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *QuizzesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *QuizzesDao) Columns() QuizzesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *QuizzesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *QuizzesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *QuizzesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
