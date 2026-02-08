package api

import "github.com/gogf/gf/v2/frame/g"

// CreateSessionReq - Create quiz session request
type CreateSessionReq struct {
	g.Meta          `path:"/sessions" tags:"Session" method:"POST" summary:"Create quiz session" description:"Create a new quiz session for a quiz." security:"BearerAuth"`
	QuizId          string `json:"quizId" v:"required#Quiz ID is required" dc:"Quiz ID to create session for"`
	MaxParticipants int    `json:"maxParticipants" v:"min:1|max:500#Max participants must be between 1 and 500" dc:"Maximum number of participants (default: 50)"`
}

type CreateSessionRes struct {
	SessionId       string `json:"sessionId" dc:"Session ID"`
	SessionCode     string `json:"sessionCode" dc:"6-character session code"`
	QuizId          string `json:"quizId" dc:"Quiz ID"`
	QuizTitle       string `json:"quizTitle" dc:"Quiz title"`
	TotalQuestions  int    `json:"totalQuestions" dc:"Total number of questions"`
	Status          string `json:"status" dc:"Session status (waiting)"`
	HostId          string `json:"hostId" dc:"Host user ID"`
	MaxParticipants int    `json:"maxParticipants" dc:"Maximum participants allowed"`
	CreatedAt       string `json:"createdAt" dc:"Creation timestamp"`
}

// JoinSessionReq - Join quiz session request
type JoinSessionReq struct {
	g.Meta      `path:"/sessions/join" tags:"Session" method:"POST" summary:"Join quiz session" description:"Join a quiz session using session code." security:"BearerAuth"`
	SessionCode string `json:"sessionCode" v:"required|length:6,6#Session code is required|Session code must be 6 characters" dc:"6-character session code"`
}

type JoinSessionRes struct {
	SessionId      string              `json:"sessionId" dc:"Session ID"`
	SessionCode    string              `json:"sessionCode" dc:"Session code"`
	QuizTitle      string              `json:"quizTitle" dc:"Quiz title"`
	Status         string              `json:"status" dc:"Session status"`
	ParticipantId  int64               `json:"participantId" dc:"Participant ID"`
	Participants   []ParticipantInfo   `json:"participants" dc:"Current participants"`
	CentrifugoInfo *CentrifugoAuthInfo `json:"centrifugoInfo" dc:"Centrifugo authentication info"`
}

type ParticipantInfo struct {
	UserId    string `json:"userId" dc:"User ID"`
	Username  string `json:"username" dc:"Username"`
	FullName  string `json:"fullName" dc:"Full name"`
	AvatarUrl string `json:"avatarUrl" dc:"Avatar URL"`
	Score     int    `json:"score" dc:"Current score"`
	Rank      int    `json:"rank" dc:"Current rank"`
	JoinedAt  string `json:"joinedAt" dc:"Join timestamp"`
}

type CentrifugoAuthInfo struct {
	Token    string   `json:"token" dc:"Centrifugo JWT token"`
	Channels []string `json:"channels" dc:"Subscribed channels"`
}

// GetSessionReq - Get session details request
type GetSessionReq struct {
	g.Meta    `path:"/sessions/{id}" tags:"Session" method:"GET" summary:"Get session details" description:"Get quiz session details." security:"BearerAuth"`
	SessionId string `json:"id" v:"required#Session ID is required" dc:"Session ID"`
}

type GetSessionRes struct {
	SessionId            string `json:"sessionId" dc:"Session ID"`
	SessionCode          string `json:"sessionCode" dc:"Session code"`
	QuizId               string `json:"quizId" dc:"Quiz ID"`
	QuizTitle            string `json:"quizTitle" dc:"Quiz title"`
	Status               string `json:"status" dc:"Session status"`
	ParticipantsCount    int    `json:"participantsCount" dc:"Number of participants"`
	CurrentQuestionIndex int    `json:"currentQuestionIndex" dc:"Current question index"`
	TotalQuestions       int    `json:"totalQuestions" dc:"Total questions"`
	StartedAt            string `json:"startedAt" dc:"Session start time"`
	EndedAt              string `json:"endedAt" dc:"Session end time"`
	DurationSeconds      int    `json:"durationSeconds" dc:"Session duration in seconds"`
	HostId               string `json:"hostId" dc:"Host user ID"`
	HostUsername         string `json:"hostUsername" dc:"Host username"`
}
