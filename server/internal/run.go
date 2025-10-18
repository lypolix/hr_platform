package internal

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hr-platform-mosprom/internal/adapter/bcrypt"
	"github.com/hr-platform-mosprom/internal/adapter/ginhandler"
	"github.com/hr-platform-mosprom/internal/adapter/jwt"
	"github.com/hr-platform-mosprom/internal/adapter/postgres"
	"github.com/hr-platform-mosprom/internal/adapter/postgres/pgqueries"
	clock "github.com/hr-platform-mosprom/internal/adapter/time"
	"github.com/hr-platform-mosprom/internal/core/application/service"
)

func Run(env environment) error {
	fmt.Println("hello world")

	db, err := postgres.New(env.PostgresDSN)
	if err != nil {
		return fmt.Errorf("error creating postgres db connection: %w", err)
	}
	defer db.Close()

	queries := pgqueries.New(db)

	utcClock := clock.NewUTCClock()
	postgresUniversityRepo := postgres.NewUniversityRepo(queries)
	jwtService := jwt.NewJWTService(utcClock, time.Hour*72, env.SecretKey)
	bcryptPasswordService := bcrypt.NewBcryptPasswordService(12)
	universityService := service.NewUniversityService(
		postgresUniversityRepo,
		bcryptPasswordService,
		jwtService,
		utcClock,
	)

	validator := validator.New()
	logger := slog.New(
		slog.Handler(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level:     slog.LevelDebug,
					AddSource: true,
				},
			),
		),
	)

	engine := gin.Default()

	ginhandler.RegisterUniversityHandlers(
		engine,
		universityService,
		logger,
		validator,
	)

	err = engine.Run(":80")
	if err != nil {
		return fmt.Errorf("error running gin engine: %w", err)
	}

	return nil
}
