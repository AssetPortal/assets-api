package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Token struct {
	bun.BaseModel `bun:"table:tokens,alias:t"`
	ID            *int      `bun:"id" json:"id"`
	Token         string    `bun:"token,unique,notnull" json:"token"`
	CreatedAt     time.Time `bun:"created_at" json:"created_at"`
	ExpiresAt     time.Time `bun:"expires_at" json:"expires_at"`
	Used          bool      `bun:"used" json:"used"`
}

func (t *Token) IsValid() bool {
	return !(time.Now().Compare(t.ExpiresAt) > 1 || t.Used)
}
