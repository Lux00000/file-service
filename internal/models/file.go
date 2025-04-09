package models

import "time"

// File ...
type File struct {
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Data      []byte
}

// FileRepository ...
type FileRepository interface {
	Save(file *File) error
	Get(name string) (*File, error)
	List() ([]*File, error)
}
