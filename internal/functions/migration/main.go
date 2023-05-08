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

// Base represents base columns for all tables
// Do not use gorm.Model because of uint ID
type Base struct {
	ID        int64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BaseWithoutID represents base columns for all tables that without ID included
// Do not use gorm.Model because of uint ID
type BaseWithoutID struct {
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
		// create initial table(s)
		{
			ID: "202305051458",
			Migrate: func(tx *gorm.DB) error {
				// it's a good pratice to copy the struct inside the function,
				// so side effects are prevented if the original struct changes during the time

				type Session struct {
					Base
					Code      string `json:"code" gorm:"type:varchar(20);unique_index"`
					IPAddress string `json:"ip_address" gorm:"type:varchar(45)" `
					UserAgent string `json:"user_agent" gorm:"type:text"`

					RefreshToken string     `json:"-" gorm:"type:varchar(100);unique_index"`
					LastLogin    *time.Time `json:"last_login"`
				}

				// Drop existing table if there is, in case re-apply this migration
				if err := tx.Migrator().DropTable(&Session{}); err != nil {
					return err
				}

				if err := tx.Set("gorm:table_options", defaultTableOpts).AutoMigrate(&Session{}); err != nil {
					return err
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("sessions")
			},
		},
	})

	return nil
}
