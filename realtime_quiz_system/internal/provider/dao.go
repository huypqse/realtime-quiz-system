package provider

import (
	"realtime_quiz_system/internal/dao"

	"go.uber.org/fx"
)

// ProvideDAOs provides all DAO instances for dependency injection
// These are pre-initialized global instances from the dao package
func ProvideDAOs() fx.Option {
	return fx.Options(
		// Provide individual DAO instances
		// Since DAOs are already initialized, we just return them
		fx.Supply(&dao.Users),
		fx.Supply(&dao.Quizzes),
		fx.Supply(&dao.Questions),
		fx.Supply(&dao.AnswerOptions),
		fx.Supply(&dao.QuizSessions),
		fx.Supply(&dao.SessionParticipants),
		fx.Supply(&dao.UserAnswers),
		fx.Supply(&dao.LeaderboardSnapshots),
	)
}
