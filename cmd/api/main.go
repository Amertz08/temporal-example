package main

import (
	"log/slog"

	"github.com/Amertz08/temporal-example/internal/database"
)

type CaseRepository interface {
	Save(database.Case) (string, error)
	Get(string) (database.Case, error)
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
