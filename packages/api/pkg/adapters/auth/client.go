package auth

import (
	"context"

	"github.com/AssetPortal/assets-api/pkg/model"
)

type Client interface {
	VerifySignature(ctx context.Context, message, address, signature string) (*model.Auth, error)
}
