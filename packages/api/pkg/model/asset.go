package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Asset struct {
	bun.BaseModel `bun:"table:assets,alias:a"`
	ID_           *int               `bun:"_id" json:"_id"`
	ID            string             `bun:"id" json:"id"`
	Address       string             `bun:"address" json:"address"`
	Blockchain    *string            `bun:"blockchain" json:"blockchain"`
	Description   *string            `bun:"description" json:"description,omitempty"`
	Image         *string            `bun:"image" json:"image,omitempty"`
	Social        *map[string]string `bun:"social" json:"social,omitempty"`
	CreatedAt     *time.Time         `bun:"created_at" json:"created_at,omitempty"`
	UpdatedAt     *time.Time         `bun:"updated_at" json:"updated_at,omitempty"`
}
