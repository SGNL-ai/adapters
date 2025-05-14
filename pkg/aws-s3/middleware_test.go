// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst

package awss3_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

const (
	emptyCSVFileCode       = 800
	headersOnlyCSVFileCode = 801
)

// Custom middleware to mock S3 responses.
type mockS3Middleware struct {
	headStatusCode   int
	getStatusCode    int
	headObjectOutput *s3.HeadObjectOutput
}

func (m *mockS3Middleware) ID() string {
	return "MockS3Middleware"
}

func (m *mockS3Middleware) HandleSerialize(
	ctx context.Context,
	in middleware.SerializeInput,
	next middleware.SerializeHandler,
) (
	out middleware.SerializeOutput,
	metadata middleware.Metadata,
	err error,
) {
	switch in.Parameters.(type) {
	case *s3.HeadObjectInput:
		return m.mockHeadObject(ctx, in, next)
	case *s3.GetObjectInput:
		return m.mockGetObject(ctx, in, next)
	default:
		return next.HandleSerialize(ctx, in)
	}
}

func (m *mockS3Middleware) mockHeadObject(
	_ context.Context,
	_ middleware.SerializeInput,
	_ middleware.SerializeHandler,
) (
	out middleware.SerializeOutput,
	metadata middleware.Metadata,
	err error,
) {
	switch m.headStatusCode {
	case http.StatusOK:
		if m.headObjectOutput != nil {
			out.Result = m.headObjectOutput
		} else {
			out.Result = &s3.HeadObjectOutput{
				ContentLength: aws.Int64(1234),
				ContentType:   aws.String("text/csv"),
				ETag:          aws.String("\"f8a7b3f9be0e4c3d2e1a0b9c8d7e6f5\""),
				LastModified:  aws.Time(time.Now()),
				VersionId:     aws.String("3/L4kqtJlcpXroDTDmJ+rmSpXd3dIbrHY"),
				Metadata: map[string]string{
					"example-metadata-key": "example-metadata-value",
				},
			}
		}
	case http.StatusMovedPermanently:
		err = &smithyhttp.ResponseError{
			Response: &smithyhttp.Response{Response: &http.Response{StatusCode: http.StatusMovedPermanently}},
			// nolint: lll
			Err: errors.New("permanent redirect: The bucket you are attempting to access must be addressed using the specified endpoint"),
		}
	case http.StatusForbidden:
		err = &smithyhttp.ResponseError{
			Response: &smithyhttp.Response{Response: &http.Response{StatusCode: http.StatusForbidden}},
			Err:      errors.New("AccessDenied: Access Denied"),
		}
	case http.StatusNotFound:
		err = &smithyhttp.ResponseError{
			Response: &smithyhttp.Response{Response: &http.Response{StatusCode: http.StatusNotFound}},
			Err:      errors.New("not found: The specified key does not exist"),
		}
	default:
		err = errors.New("unexpected status code")
	}

	return
}

func (m *mockS3Middleware) mockGetObject(
	_ context.Context,
	_ middleware.SerializeInput,
	_ middleware.SerializeHandler,
) (
	out middleware.SerializeOutput,
	metadata middleware.Metadata,
	err error,
) {
	switch m.getStatusCode {
	case emptyCSVFileCode:
		emptyCSVFile := ""
		out.Result = &s3.GetObjectOutput{
			Body:          io.NopCloser(strings.NewReader(emptyCSVFile)),
			ContentLength: aws.Int64(int64(len(emptyCSVFile))),
			ContentType:   aws.String("text/csv"),
			ETag:          aws.String("\"f8a7b3f9be0e4c3d2e1a0b9c8d7e6f5\""),
			LastModified:  aws.Time(time.Now()),
			VersionId:     aws.String("3/L4kqtJlcpXroDTDmJ+rmSpXd3dIbrHY"),
			Metadata: map[string]string{
				"example-metadata-key": "example-metadata-value",
			},
		}
	case headersOnlyCSVFileCode:
		out.Result = &s3.GetObjectOutput{
			Body:          io.NopCloser(strings.NewReader(headersOnlyCSVData)),
			ContentLength: aws.Int64(int64(len(headersOnlyCSVData))),
			ContentType:   aws.String("text/csv"),
			ETag:          aws.String("\"f8a7b3f9be0e4c3d2e1a0b9c8d7e6f5\""),
			LastModified:  aws.Time(time.Now()),
			VersionId:     aws.String("3/L4kqtJlcpXroDTDmJ+rmSpXd3dIbrHY"),
			Metadata: map[string]string{
				"example-metadata-key": "example-metadata-value",
			},
		}
	case -200:
		out.Result = &s3.GetObjectOutput{
			Body:          io.NopCloser(strings.NewReader(corruptCSVData)),
			ContentLength: aws.Int64(int64(len(corruptCSVData))),
			ContentType:   aws.String("text/csv"),
			ETag:          aws.String("\"f8a7b3f9be0e4c3d2e1a0b9c8d7e6f5\""),
			LastModified:  aws.Time(time.Now()),
			VersionId:     aws.String("3/L4kqtJlcpXroDTDmJ+rmSpXd3dIbrHY"),
			Metadata: map[string]string{
				"example-metadata-key": "example-metadata-value",
			},
		}
	case http.StatusOK:
		out.Result = &s3.GetObjectOutput{
			Body:          io.NopCloser(strings.NewReader(validCSVData)),
			ContentLength: aws.Int64(int64(len(validCSVData))),
			ContentType:   aws.String("text/csv"),
			ETag:          aws.String("\"f8a7b3f9be0e4c3d2e1a0b9c8d7e6f5\""),
			LastModified:  aws.Time(time.Now()),
			VersionId:     aws.String("3/L4kqtJlcpXroDTDmJ+rmSpXd3dIbrHY"),
			Metadata: map[string]string{
				"example-metadata-key": "example-metadata-value",
			},
		}
	case http.StatusMovedPermanently:
		err = &smithyhttp.ResponseError{
			Response: &smithyhttp.Response{Response: &http.Response{StatusCode: http.StatusMovedPermanently}},
			// nolint: lll
			Err: errors.New("permanent redirect: The bucket you are attempting to access must be addressed using the specified endpoint"),
		}
	case http.StatusForbidden:
		err = &smithyhttp.ResponseError{
			Response: &smithyhttp.Response{Response: &http.Response{StatusCode: http.StatusForbidden}},
			Err:      errors.New("access denied: Access Denied"),
		}
	case http.StatusNotFound:
		err = &smithyhttp.ResponseError{
			Response: &smithyhttp.Response{Response: &http.Response{StatusCode: http.StatusNotFound}},
			Err:      errors.New("no such key: The specified key does not exist"),
		}
	default:
		err = errors.New("unexpected status code")
	}

	return
}

func mockS3Config(headStatusCode, getStatusCode int, headObjectOutput *s3.HeadObjectOutput) *aws.Config {
	mockMiddleware := &mockS3Middleware{
		headStatusCode:   headStatusCode,
		getStatusCode:    getStatusCode,
		headObjectOutput: headObjectOutput,
	}

	// Create a custom AWS config with the mock middleware
	return &aws.Config{
		Region: "us-west-2",
		APIOptions: []func(*middleware.Stack) error{
			func(s *middleware.Stack) error {
				return s.Serialize.Add(mockMiddleware, middleware.After)
			},
		},
	}
}
