package models

import "time"

type File struct {
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Data      []byte
}

type FileRepository interface {
	Save(file *File) error
	Get(name string) (*File, error)
	List() ([]*File, error)
}
