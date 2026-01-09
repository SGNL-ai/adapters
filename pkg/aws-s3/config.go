// Copyright 2026 SGNL.ai, Inc.
package awss3

import (
	"context"
	"errors"

	"github.com/sgnl-ai/adapters/pkg/config"
)

var (
	DefaultFileType    = FileTypeCSV
	SupportedFileTypes = map[string]struct{}{FileTypeCSV: {}}
)

type Config struct {
	// Common configuration
	*config.CommonConfig

	// Region is the AWS region to query.
	Region string `json:"region"`

	// Bucket is the AWS S3 bucket containing the files with entity data.
	Bucket string `json:"bucket"`

	// Prefix is the prefix of the path containing the files with entity data.
	Prefix string `json:"prefix"`

	// FileType is the extension of the files containing the entity data.
	// This defaults to "csv".
	FileType *string `json:"fileType,omitempty"`
}

// ValidateConfig validates that a Config received in a GetPage call is valid.
func (c *Config) Validate(_ context.Context) error {
	switch {
	case c == nil:
		return errors.New("the request contains an empty configuration")
	case c.Region == "":
		return errors.New("the AWS Region is not set in the configuration")
	case c.Bucket == "":
		return errors.New("the request contains an empty AWS S3 bucket name in the configuration")
	default:
		return nil
	}
}
