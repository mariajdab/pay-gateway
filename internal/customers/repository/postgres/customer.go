package postgres

import "github.com/jackc/pgx/v5/pgxpool"

type CustomerRepo struct {
	Conn *pgxpool.Pool
}
