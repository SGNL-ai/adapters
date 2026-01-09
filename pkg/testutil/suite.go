// Copyright 2026 SGNL.ai, Inc.
package testutil

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CommonSuite struct {
	suite.Suite
}

func Run(t *testing.T, s suite.TestingSuite) {
	suite.Run(t, s)
}
