package rpPostgres

import (
	"context"
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/pkg/pkgPostgres"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/pkg/pkgRetry"
)

type RpPostgres struct {
	pg *pkgPostgres.Postgres
	rs *pkgRetry.Retry
}

func New(pg *pkgPostgres.Postgres, rs *pkgRetry.Retry) *RpPostgres {
	return &RpPostgres{
		pg: pg,
		rs: rs,
	}
}

func (rp *RpPostgres) CreateImage(ctx context.Context, originalPath string, format string) (*model.Image, error) {
	query := `
		INSERT INTO image (status, original_path, format)
		VALUES ($1, $2, $3)
		RETURNING id, status, original_path, processed_path, created_at
	`

	var img model.Image
	err := rp.pg.DB.QueryRowContext(ctx, query, model.StatusUploading, originalPath, format).Scan(
		&img.ID,
		&img.Status,
		&img.OriginalPath,
		&img.ProcessedPath,
		&img.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create image: %w", err)
	}

	return &img, nil
}

func (rp *RpPostgres) GetImage(ctx context.Context, id int) (*model.Image, error) {
	query := `SELECT id, status, original_path, processed_path, created_at FROM image WHERE id = $1`

	var img model.Image
	err := rp.pg.DB.QueryRowContext(ctx, query, id).Scan(
		&img.ID,
		&img.Status,
		&img.OriginalPath,
		&img.ProcessedPath,
		&img.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	return &img, nil
}

func (rp *RpPostgres) UpdateImageStatus(ctx context.Context, id int, status model.ImageStatus) error {
	query := `UPDATE image SET status = $1, updated_at = NOW() WHERE id = $2`

	_, err := rp.pg.DB.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update image status: %w", err)
	}

	return nil
}

func (rp *RpPostgres) UpdateImageProcessedPath(ctx context.Context, id int, processedPath string) error {
	query := `UPDATE image SET processed_path = $1, status = $2 WHERE id = $3`

	_, err := rp.pg.DB.ExecContext(ctx, query, processedPath, model.StatusCompleted, id)
	if err != nil {
		return fmt.Errorf("failed to update processed path: %w", err)
	}

	return nil
}

func (rp *RpPostgres) DeleteImage(ctx context.Context, id int) error {
	query := `DELETE FROM image WHERE id = $1`

	_, err := rp.pg.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	return nil
}
