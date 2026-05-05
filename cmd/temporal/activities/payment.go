package activities

import (
	"context"

	"github.com/Amertz08/temporal-example/internal/models"
)

func CalculateFeeAmount(ctx context.Context, vinDetails *models.VinDetails) (int64, error) {
	return 10000, nil
}
