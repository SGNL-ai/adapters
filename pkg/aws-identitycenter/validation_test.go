package awsidentitycenter_test

import (
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	adapter "github.com/sgnl-ai/adapters/pkg/aws-identitycenter"
)

func TestValidateGetPageRequest(t *testing.T) {
	tests := map[string]struct {
		request *framework.Request[adapter.Config]
		wantErr *framework.Error
	}{
		"valid_request": {
			request: &framework.Request[adapter.Config]{
				Auth: validAuthCredentials,
				Entity: framework.EntityConfig{
					ExternalId: adapter.User,
					Attributes: []*framework.AttributeConfig{
						{ExternalId: "UserId", Type: framework.AttributeTypeString, UniqueId: true},
					},
				},
				Config:   validConfig,
				PageSize: 10,
			},
			wantErr: nil,
		},
		"invalid_missing_auth": {
			request: &framework.Request[adapter.Config]{
				Entity:   framework.EntityConfig{ExternalId: adapter.User},
				Config:   validConfig,
				PageSize: 10,
			},
			wantErr: &framework.Error{
				Message: "Provided datasource auth is missing required AWS authorization credentials.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			},
		},
	}

	a := &adapter.Adapter{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := a.ValidateGetPageRequest(nil, tt.request)
			if (gotErr == nil) != (tt.wantErr == nil) {
				t.Fatalf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
			if gotErr != nil && gotErr.Message != tt.wantErr.Message {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
