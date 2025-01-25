package storage

import "context"

type Client interface {
	UploadFile(ctx context.Context, key string, file []byte, contentType string) (string, error)
}
