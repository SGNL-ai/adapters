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

// GetFileRange retrieves a specific byte range from the S3 object.
// This enables streaming large files without loading everything into memory.
func (s *S3Handler) GetFileRange(ctx context.Context, bucket string, key string, startByte, endByte int64) (*[]byte, error) {
	rangeHeader := fmt.Sprintf("bytes=%d-%d", startByte, endByte)

	response, err := s.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Range:  &rangeHeader,
	})

	if err != nil {
		httpResponseErr, parseErr := httpResponseFromError(err)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to download file range: %w", err)
		}

		switch httpResponseErr.Response.StatusCode {
		case http.StatusForbidden:
			return nil, fmt.Errorf("unable to download file range due to missing permissions")
		case http.StatusNotFound:
			return nil, fmt.Errorf("file does not exist")
		case http.StatusRequestedRangeNotSatisfiable:
			return nil, fmt.Errorf("requested byte range is not satisfiable")
		default:
			return nil, fmt.Errorf("failed to download file range: %w", err)
		}
	}

	if response == nil {
		return nil, fmt.Errorf("failed to download file range: response is nil")
	}

	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read file range body: %w", err)
	}

	return &bytes, nil
}

// GetFileSize returns the size of the file in bytes using HEAD request.
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

// GetHeaderChunk reads the first chunk of the file to extract CSV headers.
func (s *S3Handler) GetHeaderChunk(ctx context.Context, bucket string, key string) (*[]byte, error) {
	const headerChunkSize = 8196 // 8KB should be enough for most CSV headers

	return s.GetFileRange(ctx, bucket, key, 0, headerChunkSize-1)
}
