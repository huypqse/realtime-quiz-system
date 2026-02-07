// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Users is the golang structure of table users for DAO operations like Where/Data.
type Users struct {
	g.Meta       `orm:"table:users, do:true"`
	Id           any         //
	Username     any         //
	Email        any         //
	PasswordHash any         //
	FullName     any         //
	AvatarUrl    any         //
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
}

type UserRegisterReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterRes struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type UserLoginReq struct {
	UsernameOrEmail string `json:"usernameOrEmail"`
	Password        string `json:"password"`
}

type UserLoginRes struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FullName  string `json:"fullName"`
	AvatarUrl string `json:"avatarUrl"`
	Token     string `json:"token"`
}

type UserProfileRes struct {
	Id              string `json:"id"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	FullName        string `json:"fullName"`
	AvatarUrl       string `json:"avatarUrl"`
	SessionsPlayed  int    `json:"sessionsPlayed"`
	AverageScore    int    `json:"averageScore"`
	HighestScore    int    `json:"highestScore"`
	FirstPlaceCount int    `json:"firstPlaceCount"`
}

type UserUpdateProfileReq struct {
	FullName  string `json:"fullName"`
	AvatarUrl string `json:"avatarUrl"`
}
