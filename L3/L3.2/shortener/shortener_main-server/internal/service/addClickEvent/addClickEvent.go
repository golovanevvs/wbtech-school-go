package addClickEvent

import (
	"context"
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/model"
)

type ISaveClickEventRepository interface {
	SaveClickEvent(ctx context.Context, event model.Analitic) error
}

type AddClickEventService struct {
	rpSaveClickEvent ISaveClickEventRepository
}

func New(
	rpSaveClickEvent ISaveClickEventRepository,
) *AddClickEventService {
	return &AddClickEventService{
		rpSaveClickEvent: rpSaveClickEvent,
	}
}

func (sv *AddClickEventService) AddClickEvent(ctx context.Context, event model.Analitic) error {
	err := sv.rpSaveClickEvent.SaveClickEvent(ctx, event)
	if err != nil {
		return fmt.Errorf("failed to add click event: %w", err)
	}

	return nil
}
