package user

import "github.com/gogf/gf/v2/frame/g"

type UserRegisterReq struct {
	g.Meta   `path:"/register" tags:"User" method:"POST" summary:"Register a new user" description:"Creates a new user account."`
	Username string `json:"username" v:"required|length:3,50#Username is required|Username must be 3-50 characters" dc:"Username for the new account"`
	Email    string `json:"email" v:"required|email#Email is required|Invalid email format" dc:"Email address for the new account"`
	Password string `json:"password" v:"required|length:6,50#Password is required|Password must be 6-50 characters" dc:"Password for the account"`
}

type UserRegisterRes struct {
	Id       string `json:"id" dc:"User ID"`
	Username string `json:"username" dc:"Registered username"`
	Email    string `json:"email" dc:"Registered email address"`
	Token    string `json:"token" dc:"JWT authentication token"`
}

type UserLoginReq struct {
	g.Meta          `path:"/login" tags:"User" method:"POST" summary:"Login user" description:"Authenticate user and return JWT token."`
	UsernameOrEmail string `json:"usernameOrEmail" v:"required#Username or email is required" dc:"Username or email address"`
	Password        string `json:"password" v:"required#Password is required" dc:"User password"`
}

type UserLoginRes struct {
	Id        string `json:"id" dc:"User ID"`
	Username  string `json:"username" dc:"Username"`
	Email     string `json:"email" dc:"Email address"`
	FullName  string `json:"fullName" dc:"Full name"`
	AvatarUrl string `json:"avatarUrl" dc:"Avatar URL"`
	Token     string `json:"token" dc:"JWT authentication token"`
}

type UserProfileReq struct {
	g.Meta `path:"/profile" tags:"User" method:"GET" summary:"Get user profile" description:"Get authenticated user's profile with statistics." security:"BearerAuth"`
}

type UserProfileRes struct {
	Id              string `json:"id" dc:"User ID"`
	Username        string `json:"username" dc:"Username"`
	Email           string `json:"email" dc:"Email address"`
	FullName        string `json:"fullName" dc:"Full name"`
	AvatarUrl       string `json:"avatarUrl" dc:"Avatar URL"`
	SessionsPlayed  int    `json:"sessionsPlayed" dc:"Total number of quiz sessions played"`
	AverageScore    int    `json:"averageScore" dc:"Average score across all sessions"`
	HighestScore    int    `json:"highestScore" dc:"Highest score achieved"`
	FirstPlaceCount int    `json:"firstPlaceCount" dc:"Number of times ranked first"`
}

type UserUpdateProfileReq struct {
	g.Meta    `path:"/profile" tags:"User" method:"PUT" summary:"Update user profile" description:"Update authenticated user's profile (full_name and avatar_url only)." security:"BearerAuth"`
	FullName  string `json:"fullName" v:"max-length:100#Full name must be at most 100 characters" dc:"Full name"`
	AvatarUrl string `json:"avatarUrl" v:"max-length:255|url#Avatar URL must be at most 255 characters|Invalid URL format" dc:"Avatar URL"`
}

type UserUpdateProfileRes struct {
	Message string `json:"message" dc:"Success message"`
}
