package service

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/model"
)

func (sv *Service) GetImage(ctx context.Context, id int) (*model.Image, error) {
	return sv.rpMeta.GetImage(ctx, id)
}

func (s *Service) GetAllImages(ctx context.Context) ([]model.Image, error) {
	return s.rpMeta.GetAllImages(ctx)
}
