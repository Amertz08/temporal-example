package main

import (
	"log/slog"
)

type CaseRepository interface {
	Save(Case) (string, error)
	Get(string) (Case, error)
	Close() error
}

func main() {
	var repo CaseRepository
	var err error
	repo, err = NewJSONFileDB("cases.json")
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
