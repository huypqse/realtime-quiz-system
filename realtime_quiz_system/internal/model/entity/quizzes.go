// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

// Quizzes is the golang structure for table quizzes.
type Quizzes struct {
	Id               uuid.UUID   `json:"id"               orm:"id"                 description:""`
	Title            string      `json:"title"            orm:"title"              description:""`
	Description      string      `json:"description"      orm:"description"        description:""`
	CreatedBy        uuid.UUID   `json:"createdBy"        orm:"created_by"         description:""`
	Status           string      `json:"status"           orm:"status"             description:""`
	StartTime        *gtime.Time `json:"startTime"        orm:"start_time"         description:""`
	EndTime          *gtime.Time `json:"endTime"          orm:"end_time"           description:""`
	TimeLimitSeconds int         `json:"timeLimitSeconds" orm:"time_limit_seconds" description:""`
	TotalQuestions   int         `json:"totalQuestions"   orm:"total_questions"    description:""`
	CreatedAt        *gtime.Time `json:"createdAt"        orm:"created_at"         description:""`
	UpdatedAt        *gtime.Time `json:"updatedAt"        orm:"updated_at"         description:""`
}
