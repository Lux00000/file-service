package upload

import (
	"context"
	domain "file-service/internal/models"
	semaphore "file-service/internal/utils"
)

const count_downloads = 10

// UploadUseCase ...
type UploadUseCase struct {
	repo      domain.FileRepository
	semaphore *semaphore.Semaphore
}

// NewUploadUseCase ...
func NewUploadUseCase(repo domain.FileRepository) *UploadUseCase {
	return &UploadUseCase{
		repo:      repo,
		semaphore: semaphore.NewSemaphore(count_downloads),
	}
}

// Upload ...
func (uc *UploadUseCase) Upload(ctx context.Context, file *domain.File) error {
	uc.semaphore.Acquire()
	defer uc.semaphore.Release()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return uc.repo.Save(file)
	}
}
