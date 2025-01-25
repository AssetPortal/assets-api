package storage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Uploader struct {
	client     *s3.Client
	bucketName string
}

func NewS3Uploader(cfg aws.Config, bucketName string) *S3Uploader {
	client := s3.NewFromConfig(cfg)
	return &S3Uploader{
		client:     client,
		bucketName: bucketName,
	}
}

func (u *S3Uploader) UploadFile(ctx context.Context, key string, file []byte, contentType string) (string, error) {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(u.bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(file),
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPublicRead, // Optional: Adjust based on your access control requirements.
	}

	_, err := u.client.PutObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", u.bucketName, key)
	return fileURL, nil
}
