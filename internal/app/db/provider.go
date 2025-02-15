package db

import (
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ternaryss/houston/internal/app/settings"
)

var provider *dbProvider

type dbProvider struct {
	db *sql.DB
}

func NewDbProvider(stg settings.Settings) *dbProvider {
	if provider != nil {
		slog.Error("Database provider already exists")
		os.Exit(1)
	}

	file := stg.Database.File
	slog.Info("Opening SQLite connection", "file", file)
	db, err := sql.Open("sqlite3", file)

	if err != nil {
		slog.Error("Connection to SQLite failed", "err", err)
		os.Exit(1)
	}

	if _, err := db.Exec("PRAGMA FOREIGN_KEYS=ON;"); err != nil {
		slog.Error("Setting foreign keys for SQLite failed", "err", err)
		os.Exit(1)
	}

	if err := db.Ping(); err != nil {
		slog.Error("Pinging SQLite database failed", "err", err)
		os.Exit(1)
	}

	slog.Info("SQLite connection established")
	provider = &dbProvider{
		db: db,
	}

	return provider
}

func (p *dbProvider) CloseConnection() {
	if err := p.db.Close(); err != nil {
		slog.Error("SQLite connection closing failed", "err", err)
		os.Exit(1)
	}
}
