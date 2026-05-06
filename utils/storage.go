package utils

import (
	"context"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

// Storage defines operations needed by HTTP handlers.
type Storage interface {
	UploadObject(ctx context.Context, key string, body io.Reader, contentType string) error
	GetObject(ctx context.Context, key string) (*Object, error)
	IsNotFound(err error) bool
}

// Object represents a streamable file from object storage.
type Object struct {
	Body          io.ReadCloser
	ContentType   string
	ContentLength int64
	ETag          string
}

// S3Storage is AWS S3 implementation for Storage.
type S3Storage struct {
	client *s3.Client
	bucket string
}

// NewS3Storage creates S3 storage client from app config.
func NewS3Storage(cfg *Config) (*S3Storage, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion(cfg.AWSRegion),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AWSAccessKeyID,
				cfg.AWSSecretAccessKey,
				cfg.AWSSessionToken,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.AWSEndpointURL != "" {
			o.BaseEndpoint = aws.String(cfg.AWSEndpointURL)
		}
		o.UsePathStyle = cfg.AWSUsePathStyle
	})

	return &S3Storage{client: client, bucket: cfg.AWSBucket}, nil
}

// UploadObject uploads object stream to bucket.
func (s *S3Storage) UploadObject(ctx context.Context, key string, body io.Reader, contentType string) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
		ACL:         s3types.ObjectCannedACLPublicRead,
	})
	return err
}

// GetObject fetches object stream by key.
func (s *S3Storage) GetObject(ctx context.Context, key string) (*Object, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	obj := &Object{Body: out.Body}
	if out.ContentType != nil {
		obj.ContentType = *out.ContentType
	}
	if out.ContentLength != nil {
		obj.ContentLength = *out.ContentLength
	}
	if out.ETag != nil {
		obj.ETag = *out.ETag
	}

	return obj, nil
}

// IsNotFound returns true when upstream indicates missing object.
func (s *S3Storage) IsNotFound(err error) bool {
	var noSuchKey *s3types.NoSuchKey
	if errors.As(err, &noSuchKey) {
		return true
	}

	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		code := apiErr.ErrorCode()
		return code == "NoSuchKey" || code == "NotFound"
	}

	return false
}
