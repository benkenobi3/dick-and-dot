package repository

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

const GetDickSQL = `
with cte(user_id, length) as (
	select user_id, sum(length)
	from dick
	where chat_id = $1 and user_id = $2
	group by user_id
)
select d.user_id, c.length, d.updated_at 
from dick d
join cte c
on c.user_id = d.user_id
where d.chat_id = $1 and d.user_id = $2
order by d.updated_at desc
limit 1;
`

func (d *dicks) GetDick(ctx context.Context, chatID, userID int64) (Dick, error) {
	result, err := d.db.QueryContext(ctx, GetDickSQL, chatID, userID)
	if err != nil {
		return Dick{}, fmt.Errorf("cannot get dick: %w", err)
	}

	dick := Dick{
		ChatID: chatID,
	}
	for result.Next() {
		err = result.Scan(&dick.UserID, &dick.Length, &dick.UpdatedAt)
	}

	return dick, nil
}

func (d *dicks) AddDick(ctx context.Context, dick Dick) error {
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
		return fmt.Errorf("cannot add dick: %w", err)
	}
	return nil
}

const GetTopDicksByChatIDSQL = `
select user_id, sum(length) 
from dicks 
where chat_id = $1 
limit $2;
`

func (d *dicks) GetTopDicksByChatID(ctx context.Context, chatID int64, limit uint64) ([]Dick, error) {
	result, err := d.db.QueryContext(ctx, GetTopDicksByChatIDSQL, chatID, limit)
	if err != nil {
		return nil, fmt.Errorf("cannot get top dicks: %w", err)
	}

	top := make([]Dick, 0)
	for result.Next() {
		dick := Dick{}
		err = result.Scan(&dick.UserID, &dick.Length)
		if err != nil {
			return nil, fmt.Errorf("scan top dick failed: %w", err)
		}
		top = append(top, dick)
	}

	err = result.Close()
	if err != nil {
		return nil, fmt.Errorf("cannot close result after reading: %w", err)
	}

	return top, nil
}
