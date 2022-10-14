package database

import (
	"context"
	"wbintern/l0/service/model"

	"github.com/jackc/pgx/v5"
)

// Database connection
type Database struct {
	Connection *pgx.Conn
}

// Insert an order into a database
func (db *Database) InsertOrder(om model.OrderModel) error {
	_, err := db.Connection.Exec(context.Background(), `insert into orders (id, data) values ($1, $2)`, om.Uid, om.Json)
	return err
}

// Establish database connection
func Connect(connStr string) (*Database, error) {
	var database Database
	dbconn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return &database, err
	}
	database.Connection = dbconn
	return &database, err
}

// Return all orders from a database as a slice
func (db *Database) GetOrders() ([]model.OrderModel, error) {
	raw, err := db.Connection.Query(context.Background(), `select * from orders`)
	if err != nil {
		return []model.OrderModel{}, err
	}

	rows := []model.OrderModel{}
	for raw.Next() {
		row := model.OrderModel{}
		err := raw.Scan(&row.Uid, &row.Json)
		if err != nil {
			return rows, err
		}
		rows = append(rows, row)
	}
	return rows, err
}
