package database

import (
	"database/sql"
	"forum/internal/config"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDb(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.DriverDb, cfg.DsnDb)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	stmts, err := os.ReadFile(cfg.MigrationPath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(string(stmts))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InsertInitialData(cfg config.Config) error {

	db, err := sql.Open(cfg.DriverDb, cfg.DsnDb)
	if err != nil {
		return err
	}
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	stmts, err := os.ReadFile(cfg.InitDataPath)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(stmts))
	if err != nil {
		return err
	}

	return nil
}
