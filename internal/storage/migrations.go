package storage

import (
	"github.com/golang-migrate/migrate/v4"
	"myapp/internal/config"
)

func MakeMigrations(cfg *config.DatabaseConfig) error {

	m, err := migrate.New(
		cfg.MigrationsDSN,
		cfg.DatabaseDSN,
	)

	defer m.Close()

	if err != nil {
		return err
	}

	migrateErr := m.Up()
	if migrateErr != nil && migrateErr != migrate.ErrNoChange {
		return migrateErr
	}

	return nil
}
