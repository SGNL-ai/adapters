// Copyright 2025 SGNL.ai, Inc.

package awss3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// BOM (Byte Order Mark) patterns for different encodings.
var (
	UTF8BOM    = []byte{0xEF, 0xBB, 0xBF}
	UTF16LEBOM = []byte{0xFF, 0xFE}
	UTF16BEBOM = []byte{0xFE, 0xFF}
	UTF32LEBOM = []byte{0xFF, 0xFE, 0x00, 0x00}
	UTF32BEBOM = []byte{0x00, 0x00, 0xFE, 0xFF}
)

func stripBOM(data []byte) []byte {
	if len(data) == 0 {
		return data
	}

	// Check for BOMs in order of longest to shortest to avoid false matches
	bomPatterns := [][]byte{
		UTF32LEBOM, // 4 bytes
		UTF32BEBOM, // 4 bytes
		UTF8BOM,    // 3 bytes
		UTF16LEBOM, // 2 bytes
		UTF16BEBOM, // 2 bytes
	}

	for _, bom := range bomPatterns {
		if len(data) >= len(bom) && bytes.HasPrefix(data, bom) {
			return data[len(bom):]
		}
	}

	return data
}

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

// GetFile retrieves the entire object from the bucket.
// It returns a 403 error if ListBucket permission is missing.
// It returns a 404 error if the object does not exist in the path.
func (s *S3Handler) GetFile(ctx context.Context, bucket string, key string) (*[]byte, error) {
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
			return nil, fmt.Errorf("failed to convert response: %w", err)
		}

		return nil, httpResponseErr
	}

	if response == nil {
		return nil, fmt.Errorf("failed to download the file: response is nil")
	}

	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read the file body: %w", err)
	}

	// Strip BOM from the beginning of the file
	cleanedBytes := stripBOM(bytes)

	return &cleanedBytes, nil
}

// GetFileRange retrieves a specific byte range from the S3 object.
func (s *S3Handler) GetFileRange(ctx context.Context, bucket string, key string, sByte, eByte int64) (*[]byte, error) {
	rangeHeader := fmt.Sprintf("bytes=%d-%d", sByte, eByte)

	response, err := s.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Range:  &rangeHeader,
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

	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read file range body: %w", err)
	}

	// Only strip BOM if we're reading from the beginning of the file
	if sByte == 0 {
		bytes = stripBOM(bytes)
	}

	return &bytes, nil
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

func (s *S3Handler) GetHeaderChunk(ctx context.Context, bucket string, key string) (*[]byte, error) {
	const headerChunkSize = 8192 // 8KB

	headerBytes, err := s.GetFileRange(ctx, bucket, key, 0, headerChunkSize-1)
	if err != nil {
		return nil, err
	}

	return headerBytes, nil
}
