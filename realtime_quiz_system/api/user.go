package user

import "github.com/gogf/gf/v2/frame/g"

type UserRegisterReq struct {
	g.Meta   `path:"/register" tags:"User" method:"POST" summary:"Register a new user" description:"Creates a new user account."`
	Username string `json:"username" v:"required|length:3,50#Username is required|Username must be 3-50 characters" dc:"Username for the new account"`
	Email    string `json:"email" v:"required|email#Email is required|Invalid email format" dc:"Email address for the new account"`
}

type UserRegisterRes struct {
	Username string `json:"username" dc:"Registered username"`
	Email    string `json:"email" dc:"Registered email address"`
}
