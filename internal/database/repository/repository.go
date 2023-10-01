package repository

import "github.com/jmoiron/sqlx"

type dicks struct {
	db *sqlx.DB
}

type Dicks interface {
}

func NewDicks(db *sqlx.DB) Dicks {
	return dicks{
		db: db,
	}
}
