package db

import "github.com/jackc/pgx/v5/pgconn"

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
)

var ErrUniqueViolation = &pgconn.PgError{
	Code: UniqueViolation,
}
