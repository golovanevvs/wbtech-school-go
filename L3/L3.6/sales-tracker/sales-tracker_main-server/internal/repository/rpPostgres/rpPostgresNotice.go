package rpPostgres

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/model"
)

// NoticeRepository implements the notice repository interface for PostgreSQL
type NoticeRepository struct {
	db *UserRepository
}

// NewNoticeRepository creates a new instance of NoticeRepository
func NewNoticeRepository(userRp *UserRepository) *NoticeRepository {
	return &NoticeRepository{db: userRp}
}

// GetByID returns a user by ID with context
func (rp *NoticeRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	return rp.db.GetByID(id)
}
