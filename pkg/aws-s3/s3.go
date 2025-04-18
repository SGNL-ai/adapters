// Copyright 2025 SGNL.ai, Inc.

package awss3

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Implementation of EntityHandler for AWS package awss3.
type S3Handler struct {
	Client *s3.Client
}

// FileExists checks if the object exists.
// It returns a 403 error if ListBucket permission is missing.
// It returns a 404 error if the object does not exist in the path.
func (s *S3Handler) FileExists(ctx context.Context, bucket string, key string) (*s3.HeadObjectOutput, error) {
	response, err := s.Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	if err != nil {
		httpResponseErr, parseErr := httpResponseFromError(err)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to check if the file exists: %w", err)
		}

		switch httpResponseErr.Response.StatusCode {
		case http.StatusForbidden:
			return nil, fmt.Errorf("unable to check if the file exists due to missing permissions")
		case http.StatusNotFound:
			return nil, fmt.Errorf("file does not exist")
		default:
			return nil, fmt.Errorf("failed to check if the file exists: %w", httpResponseErr.Err)
		}
	}

	return response, nil
}

// GetFile retrieves the object from the bucket.
// It returns a 403 error if ListBucket permission is missing.
// It returns a 404 error if the object does not exist in the path.
func (s S3Handler) GetFile(ctx context.Context, bucket string, key string) (*[]byte, error) {
	// Check if file exists.
	// Use metadata for file type, encoding and other checks in the future.
	_, err := s.FileExists(ctx, bucket, key)
	if err != nil {
		return nil, err
	}

	response, err := s.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	if err != nil {
		httpResponseErr, parseErr := httpResponseFromError(err)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to download the file: %w", err)
		}

		switch httpResponseErr.Response.StatusCode {
		case http.StatusForbidden:
			return nil, fmt.Errorf("unable to download the file due to missing permissions")
		case http.StatusNotFound:
			return nil, fmt.Errorf("file does not exist")
		default:
			return nil, fmt.Errorf("failed to download the file: %w", httpResponseErr.Err)
		}
	}

	if response == nil {
		return nil, fmt.Errorf("failed to download the file: response is nil")
	}

	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read the file body: %w", err)
	}

	return &bytes, nil
}
