package repository

import (
    "context"
    "time"
	"clicker/internal/domain/entity"
)

type ClickRepository interface {
    SaveBatch(ctx context.Context, clicks []*entity.Click) error
    GetStats(ctx context.Context, bannerID int64, from, to time.Time) ([]*entity.Click, error)
    GetTotalClicks(ctx context.Context, bannerID int64) (int64, error)
}

type ClickUseCase interface {
    Counter(ctx context.Context, bannerID int64) (int64, error)
    Stats(ctx context.Context, bannerID int64, from, to time.Time) ([]*entity.Click, error)
}
