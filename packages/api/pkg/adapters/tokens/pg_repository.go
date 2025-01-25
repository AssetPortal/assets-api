package tokens

import (
	"context"
	"fmt"
	"strings"

	"github.com/AssetPortal/assets-api/pkg/model"
	"github.com/uptrace/bun"
)

type TokensRepository struct {
	db *bun.DB
}

func NewTokensRepository(db *bun.DB) *TokensRepository {
	return &TokensRepository{db: db}
}

func (repo *TokensRepository) CreateToken(ctx context.Context, token *model.Token) (*model.Token, error) {
	_, err := repo.db.NewInsert().Model(token).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to store token in database: %v", err)
	}

	return token, nil
}

func (repo *TokensRepository) GetToken(ctx context.Context, token string) (*model.Token, error) {
	var dbToken model.Token
	err := repo.db.NewSelect().Model(&dbToken).
		Where("token = ?", token).
		Scan(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, nil
		}
		return &dbToken, fmt.Errorf("failed to query token: %v", err)
	}

	return &dbToken, nil
}

// MarkTokenAsUsed sets a token as used to prevent reuse.
func (repo *TokensRepository) MarkTokenAsUsed(ctx context.Context, token string) error {
	_, err := repo.db.NewUpdate().Model(&model.Token{}).
		Set("used = ?", true).
		Where("token = ?", token).
		Exec(ctx)
	return err
}
