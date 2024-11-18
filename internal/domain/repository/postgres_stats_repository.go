package repository

import (
	"context"
	"database/sql"
	"clicker/internal/domain/entity"
	"fmt"
	"time"
)

type PostgresStatsRepository struct {
	db *sql.DB
}

func NewPostgresStatsRepository(db *sql.DB) StatsRepository {
	return &PostgresStatsRepository{db: db}
}

func (r *PostgresStatsRepository) GetStats(ctx context.Context, bannerID int64, from, to time.Time) ([]*entity.Click, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT banner_id, timestamp, count
		FROM clicks
		WHERE banner_id = $1 AND timestamp BETWEEN $2 AND $3
		ORDER BY timestamp ASC
	`, bannerID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to query stats: %w", err)
	}
	defer rows.Close()

	var clicks []*entity.Click
	for rows.Next() {
		var click entity.Click
		if err := rows.Scan(&click.BannerID, &click.Timestamp, &click.Count); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		clicks = append(clicks, &click)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return clicks, nil
}