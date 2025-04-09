package s3

import (
	"bytes"
	"context"
	domain "file-service/internal/models"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"os"
	"strings"
)

type S3Repository struct {
	client     *minio.Client
	bucketName string
}

func NewS3Repository(bucketName string) (*S3Repository, error) {
	endpoint := os.Getenv("S3_ENDPOINT")
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	useSSL := false

	// Remove protocol from endpoint if present
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	// Create bucket if it doesn't exist
	err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := client.BucketExists(context.Background(), bucketName)
		if errBucketExists != nil || !exists {
			return nil, err
		}
	}

	return &S3Repository{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (r *S3Repository) Save(file *domain.File) error {
	ctx := context.Background()
	reader := bytes.NewReader(file.Data)

	_, err := r.client.PutObject(ctx, r.bucketName, file.Name, reader, int64(len(file.Data)), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *S3Repository) Get(name string) (*domain.File, error) {
	ctx := context.Background()

	obj, err := r.client.GetObject(ctx, r.bucketName, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	info, err := obj.Stat()
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, err
	}

	return &domain.File{
		Name:      name,
		CreatedAt: info.LastModified,
		UpdatedAt: info.LastModified,
		Data:      data,
	}, nil
}

func (r *S3Repository) List() ([]*domain.File, error) {
	ctx := context.Background()
	objects := r.client.ListObjects(ctx, r.bucketName, minio.ListObjectsOptions{})

	var files []*domain.File
	for obj := range objects {
		if obj.Err != nil {
			return nil, obj.Err
		}

		files = append(files, &domain.File{
			Name:      obj.Key,
			CreatedAt: obj.LastModified,
			UpdatedAt: obj.LastModified,
		})
	}

	return files, nil
}
