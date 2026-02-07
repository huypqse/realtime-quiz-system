package logic

import (
	"context"
	"errors"
	"fmt"
	"math"
	"realtime_quiz_system/api"
	"realtime_quiz_system/internal/dao"
	"realtime_quiz_system/internal/model/do"
	"realtime_quiz_system/internal/model/entity"
	"realtime_quiz_system/internal/service"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/google/uuid"
)

type QuizServiceImpl struct {
	logger *glog.Logger
}

func NewQuizService(logger *glog.Logger) service.QuizService {
	return &QuizServiceImpl{
		logger: logger,
	}
}

func (s *QuizServiceImpl) CreateQuiz(ctx context.Context, userId string, req *api.CreateQuizReq) (*api.CreateQuizRes, error) {
	s.logger.Info(ctx, "Creating quiz", "userId", userId, "title", req.Title)

	// Parse user ID
	creatorId, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Validate at least one question exists
	if len(req.Questions) == 0 {
		return nil, errors.New("quiz must contain at least one question")
	}

	// Validate maximum questions limit
	if len(req.Questions) > 100 {
		return nil, errors.New("quiz cannot have more than 100 questions")
	}

	// Validate each question has exactly one correct answer
	for i, question := range req.Questions {
		// Validate answer options count (2-6)
		if len(question.AnswerOptions) < 2 {
			return nil, fmt.Errorf("question %d must have at least 2 answer options", i+1)
		}
		if len(question.AnswerOptions) > 6 {
			return nil, fmt.Errorf("question %d cannot have more than 6 answer options", i+1)
		}

		correctCount := 0
		for _, option := range question.AnswerOptions {
			if option.IsCorrect {
				correctCount++
			}
		}
		if correctCount != 1 {
			return nil, fmt.Errorf("question %d must have exactly one correct answer, found %d", i+1, correctCount)
		}

		// Set defaults
		if question.Points == 0 {
			req.Questions[i].Points = 10
		}
		if question.TimeLimitSeconds == 0 {
			req.Questions[i].TimeLimitSeconds = 30
		}
	}

	// Start transaction
	err = dao.Quizzes.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// Create quiz
		quizId, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("failed to generate quiz ID: %w", err)
		}

		_, err = dao.Quizzes.Ctx(ctx).Data(do.Quizzes{
			Id:             quizId,
			Title:          req.Title,
			Description:    req.Description,
			CreatedBy:      creatorId,
			Status:         "draft",
			TotalQuestions: len(req.Questions),
		}).Insert()
		if err != nil {
			return fmt.Errorf("failed to create quiz: %w", err)
		}

		// Create questions and answer options
		for order, question := range req.Questions {
			questionId, err := s.createQuestion(ctx, quizId, order+1, question)
			if err != nil {
				return fmt.Errorf("failed to create question %d: %w", order+1, err)
			}

			// Create answer options
			for optOrder, option := range question.AnswerOptions {
				err = s.createAnswerOption(ctx, questionId, optOrder+1, option)
				if err != nil {
					return fmt.Errorf("failed to create answer option for question %d: %w", order+1, err)
				}
			}
		}

		s.logger.Info(ctx, "Quiz created successfully", "quizId", quizId.String())
		return nil
	})

	if err != nil {
		s.logger.Error(ctx, "Failed to create quiz", "error", err)
		return nil, err
	}

	// Get the created quiz
	var createdQuiz entity.Quizzes
	err = dao.Quizzes.Ctx(ctx).Where("created_by", creatorId).OrderDesc("created_at").Limit(1).Scan(&createdQuiz)
	if err != nil {
		return nil, errors.New("failed to retrieve created quiz")
	}

	return &api.CreateQuizRes{
		Id:             createdQuiz.Id.String(),
		Title:          createdQuiz.Title,
		Description:    createdQuiz.Description,
		TotalQuestions: createdQuiz.TotalQuestions,
		Status:         createdQuiz.Status,
		CreatedBy:      createdQuiz.CreatedBy.String(),
	}, nil
}

func (s *QuizServiceImpl) createQuestion(ctx context.Context, quizId uuid.UUID, order int, req api.CreateQuestionReq) (int64, error) {
	result, err := dao.Questions.Ctx(ctx).Data(do.Questions{
		QuizId:           quizId,
		QuestionText:     req.QuestionText,
		QuestionType:     "multiple_choice",
		QuestionOrder:    order,
		Points:           req.Points,
		TimeLimitSeconds: req.TimeLimitSeconds,
		Explanation:      req.Explanation,
	}).Insert()
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *QuizServiceImpl) createAnswerOption(ctx context.Context, questionId int64, order int, option api.CreateAnswerOption) error {
	_, err := dao.AnswerOptions.Ctx(ctx).Data(do.AnswerOptions{
		QuestionId:  questionId,
		OptionText:  option.OptionText,
		IsCorrect:   option.IsCorrect,
		OptionOrder: order,
	}).Insert()
	return err
}

func (s *QuizServiceImpl) ListQuizzes(ctx context.Context, req *api.ListQuizzesReq) (*api.ListQuizzesRes, error) {
	// Set defaults
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 10
	}

	// Build query
	model := dao.Quizzes.Ctx(ctx)

	// Apply filters
	if req.Status != "" {
		model = model.Where("status", req.Status)
	}
	if req.Search != "" {
		model = model.WhereLike("title", "%"+req.Search+"%")
	}

	// Get total count
	total, err := model.Count()
	if err != nil {
		s.logger.Error(ctx, "Failed to count quizzes", "error", err)
		return nil, errors.New("failed to retrieve quizzes")
	}

	// Get paginated results
	var quizzes []entity.Quizzes
	err = model.
		OrderDesc("created_at").
		Limit(req.Size).
		Offset((req.Page - 1) * req.Size).
		Scan(&quizzes)
	if err != nil {
		s.logger.Error(ctx, "Failed to list quizzes", "error", err)
		return nil, errors.New("failed to retrieve quizzes")
	}

	// Convert to response
	items := make([]api.QuizSummary, len(quizzes))
	for i, q := range quizzes {
		items[i] = api.QuizSummary{
			Id:             q.Id.String(),
			Title:          q.Title,
			Description:    q.Description,
			TotalQuestions: q.TotalQuestions,
			Status:         q.Status,
			CreatedBy:      q.CreatedBy.String(),
			CreatedAt:      q.CreatedAt.String(),
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Size)))

	return &api.ListQuizzesRes{
		Items:      items,
		Total:      total,
		Page:       req.Page,
		Size:       req.Size,
		TotalPages: totalPages,
	}, nil
}

func (s *QuizServiceImpl) GetQuiz(ctx context.Context, quizId string) (*api.GetQuizRes, error) {
	// Parse quiz ID
	qid, err := uuid.Parse(quizId)
	if err != nil {
		return nil, errors.New("invalid quiz ID")
	}

	// Get quiz
	var quizEntity entity.Quizzes
	err = dao.Quizzes.Ctx(ctx).Where("id", qid).Scan(&quizEntity)
	if err != nil {
		s.logger.Error(ctx, "Failed to get quiz", "error", err)
		return nil, errors.New("failed to retrieve quiz")
	}
	if quizEntity.Id.String() == "00000000-0000-0000-0000-000000000000" {
		return nil, errors.New("quiz not found")
	}

	// Get questions
	var questions []entity.Questions
	err = dao.Questions.Ctx(ctx).
		Where("quiz_id", qid).
		OrderAsc("question_order").
		Scan(&questions)
	if err != nil {
		s.logger.Error(ctx, "Failed to get questions", "error", err)
		return nil, errors.New("failed to retrieve quiz questions")
	}

	// Get answer options for all questions
	questionInfos := make([]api.QuestionInfo, len(questions))
	for i, q := range questions {
		var options []entity.AnswerOptions
		err = dao.AnswerOptions.Ctx(ctx).
			Where("question_id", q.Id).
			OrderAsc("option_order").
			Scan(&options)
		if err != nil {
			s.logger.Error(ctx, "Failed to get answer options", "questionId", q.Id, "error", err)
			return nil, errors.New("failed to retrieve answer options")
		}

		optionInfos := make([]api.AnswerOptionInfo, len(options))
		for j, opt := range options {
			optionInfos[j] = api.AnswerOptionInfo{
				Id:          opt.Id,
				OptionText:  opt.OptionText,
				OptionOrder: opt.OptionOrder,
				// IsCorrect is intentionally excluded for security
			}
		}

		questionInfos[i] = api.QuestionInfo{
			Id:               q.Id,
			QuestionText:     q.QuestionText,
			QuestionOrder:    q.QuestionOrder,
			Points:           q.Points,
			TimeLimitSeconds: q.TimeLimitSeconds,
			Explanation:      q.Explanation,
			AnswerOptions:    optionInfos,
		}
	}

	return &api.GetQuizRes{
		Id:             quizEntity.Id.String(),
		Title:          quizEntity.Title,
		Description:    quizEntity.Description,
		TotalQuestions: quizEntity.TotalQuestions,
		Status:         quizEntity.Status,
		CreatedBy:      quizEntity.CreatedBy.String(),
		CreatedAt:      quizEntity.CreatedAt.String(),
		UpdatedAt:      quizEntity.UpdatedAt.String(),
		Questions:      questionInfos,
	}, nil
}

func (s *QuizServiceImpl) UpdateQuiz(ctx context.Context, userId string, req *api.UpdateQuizReq) error {
	// Parse IDs
	uid, err := uuid.Parse(userId)
	if err != nil {
		return errors.New("invalid user ID")
	}
	qid, err := uuid.Parse(req.Id)
	if err != nil {
		return errors.New("invalid quiz ID")
	}

	// Get quiz
	var quizEntity entity.Quizzes
	err = dao.Quizzes.Ctx(ctx).Where("id", qid).Scan(&quizEntity)
	if err != nil {
		return errors.New("failed to retrieve quiz")
	}
	if quizEntity.Id.String() == "00000000-0000-0000-0000-000000000000" {
		return errors.New("quiz not found")
	}

	// Authorization: only creator can update
	if quizEntity.CreatedBy != uid {
		return errors.New("unauthorized: only quiz creator can update this quiz")
	}

	// Check for active sessions
	hasActiveSessions, err := s.hasActiveSessions(ctx, qid)
	if err != nil {
		return errors.New("failed to check active sessions")
	}
	if hasActiveSessions {
		return errors.New("cannot update quiz with active sessions")
	}

	// Validate status transition if status is being changed
	if req.Status != "" && req.Status != quizEntity.Status {
		if !s.isValidStatusTransition(quizEntity.Status, req.Status) {
			return fmt.Errorf("invalid status transition from '%s' to '%s'", quizEntity.Status, req.Status)
		}
	}

	// Build update data
	updateData := do.Quizzes{}
	if req.Title != "" {
		updateData.Title = req.Title
	}
	if req.Description != "" {
		updateData.Description = req.Description
	}
	if req.Status != "" {
		updateData.Status = req.Status
	}

	// Update quiz
	_, err = dao.Quizzes.Ctx(ctx).Where("id", qid).Data(updateData).Update()
	if err != nil {
		s.logger.Error(ctx, "Failed to update quiz", "error", err)
		return errors.New("failed to update quiz")
	}

	s.logger.Info(ctx, "Quiz updated successfully", "quizId", qid.String())
	return nil
}

func (s *QuizServiceImpl) DeleteQuiz(ctx context.Context, userId string, quizId string) error {
	// Parse IDs
	uid, err := uuid.Parse(userId)
	if err != nil {
		return errors.New("invalid user ID")
	}
	qid, err := uuid.Parse(quizId)
	if err != nil {
		return errors.New("invalid quiz ID")
	}

	// Get quiz
	var quizEntity entity.Quizzes
	err = dao.Quizzes.Ctx(ctx).Where("id", qid).Scan(&quizEntity)
	if err != nil {
		return errors.New("failed to retrieve quiz")
	}
	if quizEntity.Id.String() == "00000000-0000-0000-0000-000000000000" {
		return errors.New("quiz not found")
	}

	// Authorization: only creator can delete
	if quizEntity.CreatedBy != uid {
		return errors.New("unauthorized: only quiz creator can delete this quiz")
	}

	// Check for active sessions
	hasActiveSessions, err := s.hasActiveSessions(ctx, qid)
	if err != nil {
		return errors.New("failed to check active sessions")
	}
	if hasActiveSessions {
		return errors.New("cannot delete quiz with active sessions (status: in_progress)")
	}

	// Delete quiz (cascade delete will handle related records)
	_, err = dao.Quizzes.Ctx(ctx).Where("id", qid).Delete()
	if err != nil {
		s.logger.Error(ctx, "Failed to delete quiz", "error", err)
		return errors.New("failed to delete quiz")
	}

	s.logger.Info(ctx, "Quiz deleted successfully", "quizId", qid.String())
	return nil
}

// Helper functions

func (s *QuizServiceImpl) hasActiveSessions(ctx context.Context, quizId uuid.UUID) (bool, error) {
	count, err := dao.QuizSessions.Ctx(ctx).
		Where("quiz_id", quizId).
		Where("status", "in_progress").
		Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *QuizServiceImpl) isValidStatusTransition(current, next string) bool {
	// Valid transitions: draft → active → completed → archived
	transitions := map[string][]string{
		"draft":     {"active"},
		"active":    {"completed"},
		"completed": {"archived"},
		"archived":  {}, // No transitions from archived
	}

	allowedNext, exists := transitions[current]
	if !exists {
		return false
	}

	for _, allowed := range allowedNext {
		if strings.EqualFold(allowed, next) {
			return true
		}
	}
	return false
}
