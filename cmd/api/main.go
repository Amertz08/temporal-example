package main

import (
	"log/slog"

	"github.com/Amertz08/temporal-example/internal/database"
	"github.com/Amertz08/temporal-example/internal/models"
)

type CaseRepository interface {
	Save(models_go.Case) (string, error)
	Get(string) (models_go.Case, error)
	Close() error
}

func main() {
	var repo CaseRepository
	var err error
	repo, err = database.NewJSONFileDB("cases.json")
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		return
	}
	defer repo.Close()

	e := NewServer(repo)

	if err := e.Start(":8080"); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
