package service

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.3/comment-tree/comment-tree_main-server/internal/model"
)

type iRepository interface {
	Save(ctx context.Context, comment *model.Comment) error
	LoadByID(ctx context.Context, id int) (*model.Comment, error)
	LoadChildren(ctx context.Context, parentID *int) ([]*model.Comment, error)
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, q string) ([]*model.Comment, error)
}

type Service struct {
	rp iRepository
}

func New(rp iRepository) *Service {
	return &Service{rp: rp}
}

func (s *Service) AddComment(ctx context.Context, comment *model.Comment) error {
	return s.rp.Save(ctx, comment)
}

func (s *Service) GetCommentsTree(ctx context.Context, parentID *int) ([]*model.Comment, error) {
	comments, err := s.rp.LoadChildren(ctx, parentID)
	if err != nil {
		return nil, err
	}

	for _, comment := range comments {
		children, err := s.GetCommentsTree(ctx, &comment.ID)
		if err != nil {
			return nil, err
		}
		comment.Children = children
	}

	return comments, nil
}

func (s *Service) RemoveComment(ctx context.Context, id int) error {
	return s.rp.Delete(ctx, id)
}

func (s *Service) FindComments(ctx context.Context, query string) ([]*model.Comment, error) {
	return s.rp.Search(ctx, query)
}
