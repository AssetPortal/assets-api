package tokens

import (
	"context"

	"github.com/AssetPortal/assets-api/pkg/model"
)

type Repository interface {
	CreateToken(ctx context.Context, token *model.Token) (*model.Token, error)
	GetToken(ctx context.Context, token string) (*model.Token, error)
	MarkTokenAsUsed(ctx context.Context, token string) error
}
