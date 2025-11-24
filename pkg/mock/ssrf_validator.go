// Copyright 2025 SGNL.ai, Inc.
package mock

import (
	"context"

	"github.com/sgnl-ai/adapters/pkg/validation"
)

// NoOpSSRFValidator is a validator that allows all URLs.
// This is intended for testing purposes only.
type NoOpSSRFValidator struct{}

// NewNoOpSSRFValidator creates a new NoOpSSRFValidator.
func NewNoOpSSRFValidator() validation.SSRFValidator {
	return &NoOpSSRFValidator{}
}

var _ validation.SSRFValidator = (*NoOpSSRFValidator)(nil)

// ValidateExternalURL always returns nil, allowing all URLs.
func (v *NoOpSSRFValidator) ValidateExternalURL(ctx context.Context, rawURL string) error {
	return nil
}
