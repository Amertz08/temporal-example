package main

import (
	"log/slog"

	"github.com/Amertz08/temporal-example/internal/api/handlers"
	"github.com/Amertz08/temporal-example/internal/database"
)

func main() {
	var repo handlers.CaseRepository
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
