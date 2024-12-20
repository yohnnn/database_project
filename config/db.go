package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL драйвер
)

// ConnectDB подключается к базе данных и возвращает подключение
func ConnectDB() (*sql.DB, error) {
	dsn := "host=localhost user=postgres password=9142 dbname=rzt port=5432 sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err) // ошибка с маленькой буквы
	}
	// Синхронный пинг базы данных для проверки подключения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err) // ошибка с маленькой буквы
	}
	log.Println("Database connection established")
	return db, nil
}

// Асинхронный запрос к базе данных
func QueryAsync(db *sql.DB, ctx context.Context, query string, ch chan<- *sql.Rows) {
	go func() {
		defer close(ch) // Закрываем канал по завершению
		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			log.Printf("error executing query: %v", err) // ошибка с маленькой буквы
			ch <- nil
			return
		}
		ch <- rows
	}()
}
