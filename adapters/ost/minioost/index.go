package minioost

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type St struct {
	Client *minio.Client
}

func New(uri, keyId, key string, secure bool) (*St, error) {
	client, err := minio.New(uri, &minio.Options{
		Creds:  credentials.NewStaticV4(keyId, key, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, fmt.Errorf("minioost.New error: %w", err)
	}

	return &St{
		Client: client,
	}, nil
}

func (o *St) CreateBucket(name string) error {
	ctx := context.Background()

	// check if already exists
	exists, err := o.Client.BucketExists(ctx, name)
	if err != nil {
		return fmt.Errorf("minioost.BucketExists error: %w", err)
	}
	if exists {
		return nil
	}

	err = o.Client.MakeBucket(ctx, name, minio.MakeBucketOptions{})
	if err != nil {
		return fmt.Errorf("minioost.MakeBucket error: %w", err)
	}

	return nil
}

func (o *St) PutObject(bucketName, name string, data io.Reader, contentType string) error {
	_, err := o.Client.PutObject(context.Background(), bucketName, name, data, -1, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("minioost.FPutObject error: %w", err)
	}

	return nil
}

func (o *St) GetObject(bucketName, name string) (io.ReadCloser, error) {
	result, err := o.Client.GetObject(context.Background(), bucketName, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("minioost.GetObject error: %w", err)
	}

	return result, nil
}

func (o *St) RemoveObject(bucketName, name string) error {
	err := o.Client.RemoveObject(context.Background(), bucketName, name, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("minioost.RemoveObject error: %w", err)
	}

	return nil
}
