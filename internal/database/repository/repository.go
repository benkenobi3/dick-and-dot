package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type dicks struct {
	db *sqlx.DB
}

type Dicks interface {
	AddDick(ctx context.Context, dick Dick) error
	GetDick(ctx context.Context, chatID, userID int64) (Dick, error)
	GetTopDicksByChatID(ctx context.Context, chatID int64, limit uint64) ([]Dick, error)
}

func NewDicks(db *sqlx.DB) Dicks {
	return &dicks{
		db: db,
	}
}
