package mock

import (
	"bytes"
	"io"
)

type St struct {
}

func New() *St {
	return &St{}
}

func (o *St) CreateBucket(name string) error {
	return nil
}

func (o *St) PutObject(bucketName, name string, data io.Reader, contentType string) error {
	return nil
}

func (o *St) GetObject(bucketName, name string) (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader([]byte{})), nil
}

func (o *St) RemoveObject(bucketName, name string) error {
	return nil
}
