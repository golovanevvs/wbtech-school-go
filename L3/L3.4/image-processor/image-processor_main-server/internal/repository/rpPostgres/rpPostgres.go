package rpPostgres

import (
	"context"
	"database/sql"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.3/comment-tree/comment-tree_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.3/comment-tree/comment-tree_main-server/internal/pkg/pkgPostgres"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.3/comment-tree/comment-tree_main-server/internal/pkg/pkgRetry"
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

func (rp *RpPostgres) Save(ctx context.Context, comment *model.Comment) error {
	query := `INSERT INTO comment (parent_id, text, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := rp.pg.DB.Master.QueryRowContext(ctx, query, comment.ParentID, comment.Text, comment.CreatedAt, comment.UpdatedAt).Scan(&id)
	if err != nil {
		return err
	}
	comment.ID = id
	return nil
}

func (rp *RpPostgres) LoadByID(ctx context.Context, id int) (*model.Comment, error) {
	query := `SELECT id, parent_id, text, created_at, updated_at FROM comment WHERE id = $1`
	row := rp.pg.DB.Master.QueryRowContext(ctx, query, id)
	var c model.Comment
	err := row.Scan(&c.ID, &c.ParentID, &c.Text, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (rp *RpPostgres) LoadChildren(ctx context.Context, parentID *int) ([]*model.Comment, error) {
	var query string
	var rows *sql.Rows
	var err error

	if parentID == nil {
		query = `SELECT id, parent_id, text, created_at, updated_at FROM comment WHERE parent_id IS NULL ORDER BY created_at ASC`
		rows, err = rp.pg.DB.Master.QueryContext(ctx, query)
	} else {
		query = `SELECT id, parent_id, text, created_at, updated_at FROM comment WHERE parent_id = $1 ORDER BY created_at ASC`
		rows, err = rp.pg.DB.Master.QueryContext(ctx, query, *parentID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		var c model.Comment
		err := rows.Scan(&c.ID, &c.ParentID, &c.Text, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}

	return comments, nil
}

func (rp *RpPostgres) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM comment WHERE id = $1`
	_, err := rp.pg.DB.Master.ExecContext(ctx, query, id)
	return err
}

func (rp *RpPostgres) Search(ctx context.Context, q string) ([]*model.Comment, error) {
	query := `SELECT id, parent_id, text, created_at, updated_at FROM comment WHERE text ILIKE $1`
	rows, err := rp.pg.DB.Master.QueryContext(ctx, query, "%"+q+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		var c model.Comment
		err := rows.Scan(&c.ID, &c.ParentID, &c.Text, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}

	return comments, nil
}
