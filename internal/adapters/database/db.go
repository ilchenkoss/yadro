package database

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"myapp/internal/config"
)

type DB struct {
	*sql.DB
	Cfg *config.DatabaseConfig
}

func NewConnection(cfg *config.DatabaseConfig) (*DB, error) {
	db, err := sql.Open("sqlite3", cfg.DatabasePath)

	if err != nil {
		return nil, err
	}

	return &DB{DB: db, Cfg: cfg}, nil
}

func (d *DB) Ping() error {
	err := d.DB.Ping()
	return err
}

func (d *DB) CloseConnection() error {
	err := d.DB.Close()
	return err
}

func (d *DB) MakeMigrations() error {

	m, err := migrate.New(
		d.Cfg.MigrationsDSN,
		d.Cfg.DatabaseDSN,
	)

	defer func(m *migrate.Migrate) {
		if cErr, _ := m.Close(); cErr != nil {
			//nothing
			return
		}
	}(m)

	if err != nil {
		return err
	}

	migrateErr := m.Up()
	if migrateErr != nil && migrateErr != migrate.ErrNoChange {
		return migrateErr
	}

	return nil
}
