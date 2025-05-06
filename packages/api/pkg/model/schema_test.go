package model_test

import (
	"testing"

	"github.com/AssetPortal/assets-api/pkg/model"
	"github.com/ggicci/httpin"
)

func TestCreateAssetInputValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   model.CreateAssetInput
		wantErr bool
	}{
		{
			name: "valid input",
			input: model.CreateAssetInput{
				NewAsset: model.NewAsset{
					ID:         "1a2b3c",
					Image:      strPtr("https://example.com/image.jpg"),
					Blockchain: "polkadot",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid base58 ID",
			input: model.CreateAssetInput{
				NewAsset: model.NewAsset{
					ID:    "1I0O", // Contains invalid base58 characters
					Image: strPtr("https://example.com/image.jpg"),
				},
			},
			wantErr: true,
		},
		{
			name: "invalid image URL",
			input: model.CreateAssetInput{
				NewAsset: model.NewAsset{
					ID:    "1a2b3c",
					Image: strPtr("invalid-url"),
				},
			},
			wantErr: true,
		},
		{
			name: "invalid blockchain",
			input: model.CreateAssetInput{
				NewAsset: model.NewAsset{
					ID:         "1a2b3c",
					Image:      strPtr("https://example.com/image.jpg"),
					Blockchain: "invalid",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("model.CreateAssetInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateAssetInputValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   model.UpdateAssetInput
		wantErr bool
	}{
		{
			name: "valid input",
			input: model.UpdateAssetInput{
				ID: "1a2b3c",
				UpdateAsset: model.UpdateAsset{
					Image: strPtr("https://example.com/image.jpg"),
				},
			},
			wantErr: false,
		},
		{
			name: "empty ID",
			input: model.UpdateAssetInput{
				UpdateAsset: model.UpdateAsset{
					Image: strPtr("https://example.com/image.jpg"),
				},
			},
			wantErr: true,
		},
		{
			name: "invalid image URL",
			input: model.UpdateAssetInput{
				ID: "1a2b3c",
				UpdateAsset: model.UpdateAsset{
					Image: strPtr("invalid-url"),
				},
			},
			wantErr: true,
		},
		{
			name: "invalid blockchain",
			input: model.UpdateAssetInput{
				ID: "1a2b3c",
				UpdateAsset: model.UpdateAsset{
					Blockchain: strPtr("invalid"),
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateAssetInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetAssetByIDInputValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   model.GetAssetByIDInput
		wantErr bool
	}{
		{
			name:    "valid ID",
			input:   model.GetAssetByIDInput{ID: "1a2b3c"},
			wantErr: false,
		},
		{
			name:    "empty ID",
			input:   model.GetAssetByIDInput{},
			wantErr: true,
		},
		{
			name:    "invalid base58 ID",
			input:   model.GetAssetByIDInput{ID: "1I0O"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAssetByIDInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUploadImageInputValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   model.UploadImageInput
		wantErr bool
	}{
		{
			name: "empty ID",
			input: model.UploadImageInput{
				File: &httpin.File{},
			},
			wantErr: true,
		},
		{
			name: "invalid ID format",
			input: model.UploadImageInput{
				ID:   "invalid id!",
				File: &httpin.File{},
			},
			wantErr: true,
		},
		{
			name: "nil file",
			input: model.UploadImageInput{
				ID: "valid_id",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadImageInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper function
func strPtr(s string) *string {
	return &s
}
