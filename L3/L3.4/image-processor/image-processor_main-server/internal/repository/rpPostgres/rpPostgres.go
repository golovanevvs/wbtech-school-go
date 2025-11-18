package rpPostgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/pkg/pkgErrors"
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

func (rp *RpPostgres) CreateImage(ctx context.Context, originalPath, format string, options model.ProcessOptions) (*model.Image, error) {
	opsJSON, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operations: %w", err)
	}

	query := `
		INSERT INTO image (status, original_path, format, operations)
		VALUES ($1, $2, $3, $4)
		RETURNING id, status, original_path, processed_path, created_at, operations
	`

	var img model.Image
	var processedPath *string
	var opsData []byte

	err = rp.pg.DB.QueryRowContext(ctx, query, model.StatusUploading, originalPath, format, opsJSON).Scan(
		&img.ID,
		&img.Status,
		&img.OriginalPath,
		&processedPath,
		&img.CreatedAt,
		&opsData,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create image with operations: %w", err)
	}

	img.ProcessedPath = processedPath
	err = json.Unmarshal(opsData, &img.Operations)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal operations: %w", err)
	}

	return &img, nil
}

func (rp *RpPostgres) GetImage(ctx context.Context, id int) (*model.Image, error) {
	query := `SELECT id, status, original_path, processed_path, created_at, operations FROM image WHERE id = $1`

	var img model.Image
	var processedPath *string
	var opsData []byte

	err := rp.pg.DB.QueryRowContext(ctx, query, id).Scan(
		&img.ID,
		&img.Status,
		&img.OriginalPath,
		&processedPath,
		&img.CreatedAt,
		&opsData,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("image with id %d not found: %w", id, pkgErrors.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	img.ProcessedPath = processedPath
	if opsData != nil {
		err = json.Unmarshal(opsData, &img.Operations)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal operations: %w", err)
		}
	}

	return &img, nil
}

func (rp *RpPostgres) GetAllImages(ctx context.Context) ([]model.Image, error) {
	query := `SELECT id, status, original_path, processed_path, created_at, operations FROM image ORDER BY created_at DESC`

	rows, err := rp.pg.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all images: %w", err)
	}
	defer rows.Close()

	var images []model.Image
	for rows.Next() {
		var img model.Image
		var processedPath *string
		var opsData []byte

		err := rows.Scan(
			&img.ID,
			&img.Status,
			&img.OriginalPath,
			&processedPath,
			&img.CreatedAt,
			&opsData,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan image: %w", err)
		}

		img.ProcessedPath = processedPath

		if opsData != nil {
			err = json.Unmarshal(opsData, &img.Operations)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal operations: %w", err)
			}
		}

		images = append(images, img)
	}

	return images, nil
}

func (rp *RpPostgres) UpdateImageStatus(ctx context.Context, id int, status model.ImageStatus) error {
	query := `UPDATE image SET status = $1, updated_at = NOW() WHERE id = $2`

	_, err := rp.pg.DB.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update image status: %w", err)
	}

	return nil
}

func (rp *RpPostgres) UpdateOriginalPath(ctx context.Context, id int, originalPath string) error {
	query := `UPDATE image SET original_path = $1 WHERE id = $2`

	_, err := rp.pg.DB.ExecContext(ctx, query, originalPath, id)
	if err != nil {
		return fmt.Errorf("failed to update original path: %w", err)
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
