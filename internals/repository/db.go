package repository

import (
	"database/sql"

	"github.com/vedant-colab/user-profile-app/internals/app"
	"github.com/vedant-colab/user-profile-app/internals/config"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	App *app.App
}

func InitDB() (*sql.DB, error) {
	dsn := config.Cfg.DSN
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
