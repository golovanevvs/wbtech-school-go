package repository

import (
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.3/comment-tree/comment-tree_main-server/internal/pkg/pkgPostgres"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.3/comment-tree/comment-tree_main-server/internal/pkg/pkgRetry"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.3/comment-tree/comment-tree_main-server/internal/repository/rpPostgres"
)

type Repository struct {
	*rpPostgres.RpPostgres
}

func New(pg *pkgPostgres.Postgres, rs *pkgRetry.Retry) (*Repository, error) {
	return &Repository{
		RpPostgres: rpPostgres.New(pg, rs),
	}, nil

}
