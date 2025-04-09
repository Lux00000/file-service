package download

import (
	"context"
	domain "file-service/internal/models"
	semaphore "file-service/internal/utils"
)

type DownloadUseCase struct {
	repo      domain.FileRepository
	semaphore *semaphore.Semaphore
}

func NewDownloadUseCase(repo domain.FileRepository) *DownloadUseCase {
	return &DownloadUseCase{
		repo:      repo,
		semaphore: semaphore.NewSemaphore(10), // 10 concurrent downloads
	}
}

func (uc *DownloadUseCase) Download(ctx context.Context, name string) (*domain.File, error) {
	uc.semaphore.Acquire()
	defer uc.semaphore.Release()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return uc.repo.Get(name)
	}
}
