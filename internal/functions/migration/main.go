package migration

import (
	"fmt"
	"time"

	"dullahan/config"
	dbutil "dullahan/internal/util/db"

	"github.com/M15t/ghoul/pkg/util/migration"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// EnablePostgreSQL: remove this and all tx.Set() functions bellow
// var defaultTableOpts = "ENGINE=InnoDB ROW_FORMAT=DYNAMIC"
var defaultTableOpts = ""

// Base represents base columns for all tables. Do not use gorm.Model because of uint ID
type Base struct {
	ID        int `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Run executes the migration
func Run() (respErr error) {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	db, err := dbutil.New(cfg.DbDsn, false)
	if err != nil {
		return err
	}
	// connection.Close() is not available for GORM 1.20.0
	// defer db.Close()

	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				respErr = fmt.Errorf("%s", x)
			case error:
				respErr = x
			default:
				respErr = fmt.Errorf("unknown error: %+v", x)
			}
		}
	}()

	// EnablePostgreSQL: remove these
	// workaround for "Index column size too large" error on migrations table
	initSQL := "CREATE TABLE IF NOT EXISTS migrations (id VARCHAR(255) PRIMARY KEY) " + defaultTableOpts
	if err := db.Exec(initSQL).Error; err != nil {
		return err
	}

	migration.Run(db, []*gormigrate.Migration{
		// create initial tables
		{
			ID: "201905051012",
			Migrate: func(tx *gorm.DB) error {
				// it's a good pratice to copy the struct inside the function,
				// so side effects are prevented if the original struct changes during the time

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("")
			},
		},
	})

	return nil
}
