package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type dicks struct {
	db *sqlx.DB
}

type Dicks interface {
	CreateDick(ctx context.Context, dick Dick) error
	UpdateDick(ctx context.Context, dick Dick) error
	GetDicksByChatId(ctx context.Context, chatID int64) (map[int64]Dick, error)
}

func NewDicks(db *sqlx.DB) Dicks {
	return &dicks{
		db: db,
	}
}
