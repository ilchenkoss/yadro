package database

import (
	"database/sql"
	"embed"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
	"myapp/internal-xkcd/config"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

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

	srcDriver, drErr := iofs.New(migrationsFS, "migrations")
	if drErr != nil {
		return drErr
	}

	m, err := migrate.NewWithSourceInstance("iofs", srcDriver, d.Cfg.DatabaseDSN)

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
