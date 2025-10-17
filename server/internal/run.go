package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/hr-platform-mosprom/internal/adapter/postgres"
	"github.com/hr-platform-mosprom/internal/adapter/postgres/pgqueries"
	"github.com/hr-platform-mosprom/internal/core/domain"
)

func Run(env environment) error {
	fmt.Println("hello world")

	db, err := postgres.New(env.PostgresDSN)
	if err != nil {
		return fmt.Errorf("error creating postgres db connection: %w", err)
	}
	defer db.Close()

	queries := pgqueries.New(db)

	universityRepo := postgres.NewUniversityRepo(queries)

	university, err := domain.CreateUniversity(domain.CreateUniversityAttrs{
		Title:        "hello world",
		Login:        "hello",
		PasswordHash: "123131321331",
		INN:          "2132131231231",
	}, time.Now())
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	universityRepo.Save(context.Background(), university)

	return nil
}
