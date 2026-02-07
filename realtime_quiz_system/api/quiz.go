package api

import "github.com/gogf/gf/v2/frame/g"

// CreateQuizReq - Quiz creation request
type CreateQuizReq struct {
	g.Meta      `path:"/quizzes" tags:"Quiz" method:"POST" summary:"Create a new quiz" description:"Create a new quiz with questions and answer options." security:"BearerAuth"`
	Title       string              `json:"title" v:"required|length:3,200#Title is required|Title must be 3-200 characters" dc:"Quiz title"`
	Description string              `json:"description" v:"max-length:1000#Description must be at most 1000 characters" dc:"Quiz description"`
	Questions   []CreateQuestionReq `json:"questions" v:"required#At least one question is required" dc:"Quiz questions"`
}

type CreateQuestionReq struct {
	QuestionText     string               `json:"questionText" v:"required|length:3,1000#Question text is required|Question text must be 3-1000 characters" dc:"Question text"`
	Points           int                  `json:"points" v:"min:1|max:1000#Points must be at least 1|Points must be at most 1000" dc:"Points for correct answer (default: 10)"`
	TimeLimitSeconds int                  `json:"timeLimitSeconds" v:"min:5|max:300#Time limit must be at least 5 seconds|Time limit must be at most 300 seconds" dc:"Time limit in seconds (default: 30)"`
	Explanation      string               `json:"explanation" v:"max-length:500#Explanation must be at most 500 characters" dc:"Explanation for the answer"`
	AnswerOptions    []CreateAnswerOption `json:"answerOptions" v:"required#Answer options are required" dc:"Answer options (2-6 options)"`
}

type CreateAnswerOption struct {
	OptionText string `json:"optionText" v:"required|length:1,500#Option text is required|Option text must be 1-500 characters" dc:"Answer option text"`
	IsCorrect  bool   `json:"isCorrect" dc:"Whether this option is correct"`
}

type CreateQuizRes struct {
	Id             string `json:"id" dc:"Quiz ID"`
	Title          string `json:"title" dc:"Quiz title"`
	Description    string `json:"description" dc:"Quiz description"`
	TotalQuestions int    `json:"totalQuestions" dc:"Total number of questions"`
	Status         string `json:"status" dc:"Quiz status"`
	CreatedBy      string `json:"createdBy" dc:"Creator user ID"`
}

// ListQuizzesReq - List quizzes with pagination and filters
type ListQuizzesReq struct {
	g.Meta `path:"/quizzes" tags:"Quiz" method:"GET" summary:"List quizzes" description:"List all quizzes with pagination and filters."`
	Page   int    `json:"page" v:"min:1#Page must be at least 1" dc:"Page number (default: 1)"`
	Size   int    `json:"size" v:"min:1|max:100#Size must be at least 1|Size must be at most 100" dc:"Page size (default: 10)"`
	Status string `json:"status" v:"in:draft,active,completed,archived#Status must be one of: draft, active, completed, archived" dc:"Filter by status"`
	Search string `json:"search" v:"max-length:100#Search term must be at most 100 characters" dc:"Search by title (case-insensitive, partial match)"`
}

type ListQuizzesRes struct {
	Items      []QuizSummary `json:"items" dc:"List of quizzes"`
	Total      int           `json:"total" dc:"Total number of quizzes"`
	Page       int           `json:"page" dc:"Current page"`
	Size       int           `json:"size" dc:"Page size"`
	TotalPages int           `json:"totalPages" dc:"Total number of pages"`
}

type QuizSummary struct {
	Id             string `json:"id" dc:"Quiz ID"`
	Title          string `json:"title" dc:"Quiz title"`
	Description    string `json:"description" dc:"Quiz description"`
	TotalQuestions int    `json:"totalQuestions" dc:"Total number of questions"`
	Status         string `json:"status" dc:"Quiz status"`
	CreatedBy      string `json:"createdBy" dc:"Creator user ID"`
	CreatedAt      string `json:"createdAt" dc:"Creation timestamp"`
}

// GetQuizReq - Get quiz details
type GetQuizReq struct {
	g.Meta `path:"/quizzes/{id}" tags:"Quiz" method:"GET" summary:"Get quiz details" description:"Get full quiz details including questions and options (is_correct flag hidden)."`
	Id     string `json:"id" v:"required#Quiz ID is required" dc:"Quiz ID"`
}

type GetQuizRes struct {
	Id             string         `json:"id" dc:"Quiz ID"`
	Title          string         `json:"title" dc:"Quiz title"`
	Description    string         `json:"description" dc:"Quiz description"`
	TotalQuestions int            `json:"totalQuestions" dc:"Total number of questions"`
	Status         string         `json:"status" dc:"Quiz status"`
	CreatedBy      string         `json:"createdBy" dc:"Creator user ID"`
	CreatedAt      string         `json:"createdAt" dc:"Creation timestamp"`
	UpdatedAt      string         `json:"updatedAt" dc:"Last update timestamp"`
	Questions      []QuestionInfo `json:"questions" dc:"Quiz questions"`
}

type QuestionInfo struct {
	Id               int64              `json:"id" dc:"Question ID"`
	QuestionText     string             `json:"questionText" dc:"Question text"`
	QuestionOrder    int                `json:"questionOrder" dc:"Question order"`
	Points           int                `json:"points" dc:"Points for correct answer"`
	TimeLimitSeconds int                `json:"timeLimitSeconds" dc:"Time limit in seconds"`
	Explanation      string             `json:"explanation" dc:"Explanation for the answer"`
	AnswerOptions    []AnswerOptionInfo `json:"answerOptions" dc:"Answer options"`
}

type AnswerOptionInfo struct {
	Id          int64  `json:"id" dc:"Option ID"`
	OptionText  string `json:"optionText" dc:"Answer option text"`
	OptionOrder int    `json:"optionOrder" dc:"Option order"`
	// Note: IsCorrect is intentionally excluded for security
}

// UpdateQuizReq - Update quiz
type UpdateQuizReq struct {
	g.Meta      `path:"/quizzes/{id}" tags:"Quiz" method:"PUT" summary:"Update quiz" description:"Update quiz title, description, or status (only creator can update)." security:"BearerAuth"`
	Id          string `json:"id" v:"required#Quiz ID is required" dc:"Quiz ID"`
	Title       string `json:"title" v:"length:3,200#Title must be 3-200 characters" dc:"Quiz title"`
	Description string `json:"description" v:"max-length:1000#Description must be at most 1000 characters" dc:"Quiz description"`
	Status      string `json:"status" v:"in:draft,active,completed,archived#Status must be one of: draft, active, completed, archived" dc:"Quiz status"`
}

type UpdateQuizRes struct {
	Message string `json:"message" dc:"Success message"`
}

// DeleteQuizReq - Delete quiz
type DeleteQuizReq struct {
	g.Meta `path:"/quizzes/{id}" tags:"Quiz" method:"DELETE" summary:"Delete quiz" description:"Delete quiz (only creator can delete, cannot delete if active sessions exist)." security:"BearerAuth"`
	Id     string `json:"id" v:"required#Quiz ID is required" dc:"Quiz ID"`
}

type DeleteQuizRes struct {
	Message string `json:"message" dc:"Success message"`
}
