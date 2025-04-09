package list

import (
	"context"
	domain "file-service/internal/models"
	semaphore "file-service/internal/utils"
)

type ListUseCase struct {
	repo      domain.FileRepository
	semaphore *semaphore.Semaphore
}

func NewListUseCase(repo domain.FileRepository) *ListUseCase {
	return &ListUseCase{
		repo:      repo,
		semaphore: semaphore.NewSemaphore(100), // 100 concurrent list requests
	}
}

func (uc *ListUseCase) List(ctx context.Context) ([]*domain.File, error) {
	uc.semaphore.Acquire()
	defer uc.semaphore.Release()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return uc.repo.List()
	}
}
