package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	Connection *pgx.Conn
}

func Connect(connStr string) (*Database, error) {
	var database Database
	dbconn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return &database, err
	}
	database.Connection = dbconn
	return &database, err
}
