package disk

import (
	"os"
	"path/filepath"

	domain "file-service/internal/models"
)

type DiskRepository struct {
	storagePath string
}

func NewDiskRepository(storagePath string) *DiskRepository {
	return &DiskRepository{storagePath: storagePath}
}

func (r *DiskRepository) Save(file *domain.File) error {
	path := filepath.Join(r.storagePath, file.Name)
	err := os.WriteFile(path, file.Data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (r *DiskRepository) Get(name string) (*domain.File, error) {
	path := filepath.Join(r.storagePath, name)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	return &domain.File{
		Name:      name,
		CreatedAt: info.ModTime(),
		UpdatedAt: info.ModTime(),
		Data:      data,
	}, nil
}

func (r *DiskRepository) List() ([]*domain.File, error) {
	entries, err := os.ReadDir(r.storagePath)
	if err != nil {
		return nil, err
	}

	var files []*domain.File
	for _, entry := range entries {
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			files = append(files, &domain.File{
				Name:      entry.Name(),
				CreatedAt: info.ModTime(),
				UpdatedAt: info.ModTime(),
			})
		}
	}

	return files, nil
}
