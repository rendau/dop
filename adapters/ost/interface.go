package ost

import (
	"io"
)

type Ost interface {
	CreateBucket(name string) error
	PutObject(bucketName, name string, data io.Reader, contentType string) error
	GetObject(bucketName, name string) (io.ReadCloser, error)
	RemoveObject(bucketName, name string) error
}
