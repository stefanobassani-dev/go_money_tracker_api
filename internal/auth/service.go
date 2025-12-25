package auth

import "github.com/jackc/pgx/v5"

type Service struct {
	db *pgx.Conn
}

func NewService(db *pgx.Conn) *Service {
	return &Service{
		db: db,
	}
}
