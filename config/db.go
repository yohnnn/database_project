package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	dsn := "host=localhost user=postgres password=9142 dbname=rzt port=5432 sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err) // ошибка с маленькой буквы
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err) // ошибка с маленькой буквы
	}
	log.Println("Database connection established")
	return db, nil
}

func QueryAsync(db *sql.DB, ctx context.Context, query string, ch chan<- *sql.Rows) {
	go func() {
		defer close(ch)
		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			log.Printf("error executing query: %v", err)
			ch <- nil
			return
		}
		ch <- rows
	}()
}
