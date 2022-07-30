package internal

import (
	"database/sql"
	"fmt"

	migrate "github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func dbConnect(options Options) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		options.DatabaseHost,
		options.DatabaseUser,
		options.DatabasePass,
		options.DatabaseName,
		options.DatabasePort,
		options.DatabaseSSLMode,
		options.DatabaseTimezone,
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func dbAutoMigration(sqlDB *sql.DB) error {
	driver, err := migratepg.WithInstance(sqlDB, &migratepg.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file:///migrations", "postgres", driver)
	if err != nil {
		return err
	}

	m.Up()

	return nil
}
