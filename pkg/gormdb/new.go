package gormdb

import (
	"github.com/pkg/errors"
	"github.com/spongeprojects/kubebigbrother/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path"
)

func New(dialect, dsn string) (*gorm.DB, error) {
	dbi, err := newDB(dialect, dsn)
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

func newDB(dialect, dsn string) (*gorm.DB, error) {
	config := &gorm.Config{}
	switch dialect {
	case "mysql":
		return gorm.Open(mysql.Open(dsn), config)
	case "postgres":
		return gorm.Open(postgres.Open(dsn), config)
	case "sqlite":
		// dsn == "" means in-memory
		if dsn != "" {
			_ = os.MkdirAll(path.Dir(dsn), 0755)
		}
		return gorm.Open(sqlite.Open(dsn), config)
	default:
		return nil, errors.Errorf("unsupported dialect: %s", dialect)
	}
}
