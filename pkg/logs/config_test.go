// Copyright 2025 SGNL.ai, Inc.
// nolint:lll
package logs_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/sgnl-ai/adapters/pkg/logs"
)

func TestLoadConfiguration(t *testing.T) {
	tests := map[string]struct {
		inputEnvVariables map[string]string
		wantConfiguration *logs.Config
		wantError         error
	}{
		"empty_config": {
			inputEnvVariables: map[string]string{},
			wantConfiguration: &logs.Config{
				// Defaults.
				Mode:           []string{"console"},
				Level:          "INFO",
				FilePath:       "/var/log/sgnl/unconfigured.log",
				FileMaxSize:    100,
				FileMaxBackups: 10,
				FileMaxDays:    7,
			},
			wantError: nil,
		},
		"set_config": {
			inputEnvVariables: map[string]string{
				"SGNL_LOG_LEVEL":            "DEBUG",
				"SGNL_LOG_MODE":             "file,console",
				"SGNL_LOG_FILE_PATH":        "/var/log/sgnl/adapter-sgnl.log",
				"SGNL_LOG_FILE_MAX_SIZE":    "200",
				"SGNL_LOG_FILE_MAX_BACKUPS": "20",
				"SGNL_LOG_FILE_MAX_DAYS":    "14",
			},
			wantConfiguration: &logs.Config{
				Mode:           []string{"file", "console"},
				Level:          "DEBUG",
				FilePath:       "/var/log/sgnl/adapter-sgnl.log",
				FileMaxSize:    200,
				FileMaxBackups: 20,
				FileMaxDays:    14,
			},
			wantError: nil,
		},
		"invalid_config": {
			inputEnvVariables: map[string]string{
				"SGNL_LOG_FILE_MAX_SIZE": "not_a_number", // Should be an int.
			},
			wantConfiguration: nil,
			wantError:         errors.New("decoding failed due to the following error(s):\n\n'file_max_size' cannot parse value as 'int': strconv.ParseInt: parsing \"not_a_number\": invalid syntax"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			for key, value := range test.inputEnvVariables {
				t.Setenv(key, value)
			}

			gotConfiguration, gotError := logs.LoadConfig()

			if test.wantError != nil {
				if gotError == nil {
					t.Errorf("got error = nil, want '%v'", test.wantError)

					return
				}

				if gotError.Error() != test.wantError.Error() {
					t.Errorf("got error = '%v', want '%v'", gotError, test.wantError)
				}
			} else if gotError != nil {
				t.Errorf("got error = '%v', want nil", gotError)
			}

			if !reflect.DeepEqual(gotConfiguration, test.wantConfiguration) {
				t.Fatalf("got %+v, want %+v", gotConfiguration, test.wantConfiguration)
			}
		})
	}
}
