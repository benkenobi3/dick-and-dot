package repository

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

func (d *dicks) GetDicksByChatId(ctx context.Context, chatID int64) (map[int64]Dick, error) {
	result, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select(
			"user_id",
			"chat_id",
			"length",
			"updated_at",
		).
		From(tableName).
		Where(sq.Eq{
			"chat_id": chatID,
		}).
		RunWith(d.db).
		QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get dicks from database: %w", err)
	}

	userDicks := make(map[int64]Dick, 0)
	for result.Next() {
		dick := Dick{}
		err = result.Scan(&dick.UserID, &dick.ChatID, &dick.Length, &dick.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("cannot scan dick: %w", err)
		}
		userDicks[dick.UserID] = dick
	}

	err = result.Close()
	if err != nil {
		return nil, fmt.Errorf("cannot close result after reading: %w", err)
	}

	return userDicks, nil
}

func (d *dicks) CreateDick(ctx context.Context, dick Dick) error {
	_, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Insert(tableName).
		Columns(
			"user_id",
			"chat_id",
			"length",
			"updated_at",
		).
		Values(
			dick.UserID,
			dick.ChatID,
			dick.Length,
			time.Now().UTC(),
		).
		RunWith(d.db).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("cannot update dick: %w", err)
	}
	return nil
}

func (d *dicks) UpdateDick(ctx context.Context, dick Dick) error {
	_, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Update(tableName).
		Set("length", dick.Length).
		Set("updated_at", time.Now().UTC()).
		Where(sq.Eq{
			"user_id": dick.UserID,
			"chat_id": dick.ChatID,
		}).
		RunWith(d.db).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("cannot update dick: %w", err)
	}
	return nil
}
