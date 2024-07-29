package repository

import (
	"database/sql"
	"time"
)

type CommonFunction interface {
	thisDB() *sql.DB
	thisTable() string
	queryTimeout() time.Duration
	execTimeout() time.Duration
}
