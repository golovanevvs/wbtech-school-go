package service

import "context"

func (sv *Service) DeleteImage(ctx context.Context, id int) error {
	img, err := sv.rpMeta.GetImage(ctx, id)
	if err != nil {
		return err
	}

	_ = sv.rpFile.DeleteOriginal(img.OriginalPath)

	if img.ProcessedPath != nil {
		_ = sv.rpFile.DeleteProcessed(*img.ProcessedPath)
	}

	return sv.rpMeta.DeleteImage(ctx, id)
}
