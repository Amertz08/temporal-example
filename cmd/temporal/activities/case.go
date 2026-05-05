package activities

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Amertz08/temporal-example/internal/database"
	"github.com/Amertz08/temporal-example/internal/models"
)

func GetCaseById(ctx context.Context, caseId string) (*models.Case, error) {
	repo, err := database.NewJSONFileDB("../api/cases.json")
	if err != nil {
		log.Println("Failed to initialize database", err)
		return nil, errors.New("could not initialize connection to db")
	}
	defer repo.Close()

	caseRecord, err := repo.Get(caseId)
	if err != nil {
		log.Println("Failed to get case record", err)
		return nil, errors.New(fmt.Sprintf("no case for id: %s", caseRecord))
	}
	time.Sleep(100 * time.Millisecond)

	return &caseRecord, nil
}
