// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

// Users is the golang structure for table users.
type Users struct {
	Id           uuid.UUID   `json:"id"           orm:"id"            description:""`
	Username     string      `json:"username"     orm:"username"      description:""`
	Email        string      `json:"email"        orm:"email"         description:""`
	PasswordHash string      `json:"-"            orm:"password_hash" description:""`
	FullName     string      `json:"fullName"     orm:"full_name"     description:""`
	AvatarUrl    string      `json:"avatarUrl"    orm:"avatar_url"    description:""`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:""`
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:""`
}
