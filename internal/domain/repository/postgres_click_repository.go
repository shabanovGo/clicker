package repository

import (
	"context"
	"database/sql"
	"clicker/internal/domain/entity"
	_ "github.com/lib/pq"
	"fmt"
	"time"
)

type PostgresClickRepository struct {
	db *sql.DB
}

func NewPostgresClickRepository(db *sql.DB) ClickRepository {
	return &PostgresClickRepository{db: db}
}

func (r *PostgresClickRepository) SaveBatch(ctx context.Context, clicks []*entity.Click) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO clicks (banner_id, timestamp)
		VALUES ($1, $2)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, click := range clicks {
		_, err := stmt.ExecContext(ctx, click.BannerID, click.Timestamp)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *PostgresClickRepository) GetStats(ctx context.Context, bannerID int64, from, to time.Time) ([]*entity.Click, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT banner_id, timestamp
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
		if err := rows.Scan(&click.BannerID, &click.Timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		clicks = append(clicks, &click)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return clicks, nil
}

func (r *PostgresClickRepository) GetTotalClicks(ctx context.Context, bannerID int64) (int64, error) {
	var totalClicks int64
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM clicks WHERE banner_id = $1
	`, bannerID).Scan(&totalClicks)
	if err != nil {
		return 0, fmt.Errorf("failed to get total clicks: %w", err)
	}
	return totalClicks, nil
}
