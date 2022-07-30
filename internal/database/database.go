package database

import (
	"database/sql"
	"fmt"

	migrate "github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/maetad/baroness-api/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseInterface interface {
	Create(value interface{}) (tx *gorm.DB)
	First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Find(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Save(value interface{}) (tx *gorm.DB)
	Delete(value interface{}, conds ...interface{}) (tx *gorm.DB)
}

func Connect(options config.Options) (*gorm.DB, error) {
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

func AutoMigration(sqlDB *sql.DB) error {
	driver, err := migratepg.WithInstance(sqlDB, &migratepg.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return err
	}

	m.Up()

	return nil
}
