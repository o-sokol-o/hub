package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func NewDB(ctx context.Context, host, port, userName, password, dbName, sslMode string) (*sqlx.DB, error) {

	// Параметры подключения к БД
	cfg_psql := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		host, port, userName, dbName, password, sslMode)

	// reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	// defer cancel()

	db, err := sqlx.Open("postgres", cfg_psql)
	if err != nil {
		return nil, err
	}

	// Проверим подключение к БД и вернём ошибку
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
