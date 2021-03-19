package gormdb

import (
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New(dialect, url string) (*gorm.DB, error) {
	dbi, err := newDB(dialect, url)
	if err != nil {
		return nil, err
	}
	if err := dbi.AutoMigrate(
		&models.Event{},
	); err != nil {
		return nil, errors.Wrap(err, "auto migrate error")
	}
	return dbi, nil
}

func newDB(dialect, url string) (*gorm.DB, error) {
	switch dialect {
	case "mysql":
		return gorm.Open(mysql.Open(url), &gorm.Config{})
	case "postgres":
		return gorm.Open(postgres.Open(url), &gorm.Config{})
	case "sqlite":
		return gorm.Open(sqlite.Open(url), &gorm.Config{})
	default:
		return nil, errors.Errorf("unknown dialect: %s", dialect)
	}
}
