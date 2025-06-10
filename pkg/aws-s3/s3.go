// Copyright 2025 SGNL.ai, Inc.

package awss3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Handler struct {
	Client *s3.Client
}

func (s *S3Handler) FileExists(ctx context.Context, bucket string, key string) (*s3.HeadObjectOutput, error) {
	response, err := s.Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	if err != nil {
		httpResponseErr, parseErr := httpResponseFromError(err)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to convert response: %w", err)
		}

		return nil, httpResponseErr
	}

	if response == nil {
		return nil, fmt.Errorf("failed to download file range: response is nil")
	}

	return response, nil
}

func (s *S3Handler) GetObjectStream(ctx context.Context, bucket string, key string, rangeHeader *string) (
	*s3.GetObjectOutput, error) {
	getObjectInput := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
	if rangeHeader != nil && *rangeHeader != "" {
		getObjectInput.Range = rangeHeader
	}

	response, err := s.Client.GetObject(ctx, getObjectInput)
	if err != nil {
		httpResponseErr, parseErr := httpResponseFromError(err)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to convert response: %w", err)
		}

		return nil, httpResponseErr
	}

	if response == nil {
		return nil, fmt.Errorf("failed to download file range: response is nil")
	}

	return response, nil
}

func (s *S3Handler) GetFileSize(ctx context.Context, bucket string, key string) (int64, error) {
	response, err := s.FileExists(ctx, bucket, key)
	if err != nil {
		return 0, err
	}

	if response.ContentLength == nil {
		return 0, fmt.Errorf("unable to determine file size")
	}

	return *response.ContentLength, nil
}
