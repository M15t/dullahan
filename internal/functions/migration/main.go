package migration

import (
	"fmt"
	"strings"
	"time"

	"dullahan/config"
	dbutil "dullahan/internal/util/db"

	"github.com/M15t/ghoul/pkg/util/migration"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/datatypes"
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

					TotalAllIncome           float64 `json:"total_all_income"`
					TotalEssentialExpense    float64 `json:"total_essential_expense"`
					TotalNonEssentialExpense float64 `json:"total_non_essential_expense"`
					MonthlyPaymentDebt       float64 `json:"monthly_payment_debt"`
					MonthlyNetFlow           float64 `json:"monthly_net_flow"`

					CurrentBalance float64 `json:"current_balance"`

					ActualEmergencyFund   float64 `json:"actual_emergency_fund"`
					ExpectedEmergencyFund float64 `json:"expected_emergency_fund"`
					ActualRainydayFund    float64 `json:"actual_rainyday_fund"`
					ExpectedRainydayFund  float64 `json:"expected_rainyday_fund"`
					FunFund               float64 `json:"fun_fund"`
					Investment            float64 `json:"investment"`
					RetirementPlan        float64 `json:"retirement_plan"`

					IsAchivedEmergencyFund  bool `json:"is_achived_emergency_fund" gorm:"default:false"`
					IsAchivedRainydayFund   bool `json:"is_achived_rainyday_fund" gorm:"default:false"`
					IsAchivedInvestment     bool `json:"is_achived_investment" gorm:"default:false"`
					IsAchivedRetirementPlan bool `json:"is_achived_retirement_plan" gorm:"default:false"`

					Status string `json:"status" gorm:"type:varchar(10)"`
				}

				type Income struct {
					Base
					SessionID int64   `json:"session_id"`
					Amount    float64 `json:"amount"`
					Name      string  `json:"name" gorm:"type:varchar(100)"`
					Type      string  `json:"type" gorm:"type:varchar(10);default:MONTHLY"` // MONTHLY, PASSIVE
				}

				type Expense struct {
					Base
					SessionID int64   `json:"session_id"`
					Amount    float64 `json:"amount"`
					Name      string  `json:"name" gorm:"type:varchar(100)"`
					Type      string  `json:"type" gorm:"type:varchar(15);default:ESSENTIAL"` // ESSENTIAL, NON_ESSENTIAL
				}

				type Debt struct {
					Base
					SessionID       int64          `json:"session_id"`
					Name            string         `json:"name" gorm:"type:varchar(50)"`
					RemainingAmount float64        `json:"remaining_amount"`
					MonthlyPayment  float64        `json:"monthly_payment"`
					AnnualInterest  float64        `json:"annual_interest"`
					Type            string         `json:"type" gorm:"type:varchar(10);default:FIXED"` // FIXED, FIXED_AMORTIZED, FLOAT, FLOAT_AMORTIZED
					PaymentDeadline datatypes.Date `json:"payment_deadline"`
				}

				// Drop existing table if there is, in case re-apply this migration
				if err := tx.Migrator().DropTable(&Session{}, &Income{}, &Expense{}, &Debt{}); err != nil {
					return err
				}

				if err := tx.Set("gorm:table_options", defaultTableOpts).AutoMigrate(&Session{}, &Income{}, &Expense{}, &Debt{}); err != nil {
					return err
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("sessions", "incomes", "expenses", "debts")
			},
		},
		// add column "total_all_expense" to sessions table
		{
			ID: "202305101555",
			Migrate: func(tx *gorm.DB) error {
				return tx.Exec(`ALTER TABLE sessions ADD COLUMN total_all_expense DOUBLE PRECISION DEFAULT 0;`).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Exec(`ALTER TABLE sessions DROP COLUMN total_all_expense;`).Error
			},
		},
		// add column "total_monthly_payment_debt" to sessions table
		{
			ID: "202305111036",
			Migrate: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE sessions ADD COLUMN total_monthly_payment_debt DOUBLE PRECISION DEFAULT 0;`,
					`ALTER TABLE sessions DROP COLUMN monthly_payment_debt;`,
				}

				return migration.ExecMultiple(tx, strings.Join(changes, " "))
			},
			Rollback: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE sessions ADD COLUMN monthly_payment_debt DOUBLE PRECISION DEFAULT 0;`,
					`ALTER TABLE sessions DROP COLUMN total_monthly_payment_debt;`,
				}

				return migration.ExecMultiple(tx, strings.Join(changes, " "))
			},
		},
		// add column "forecast_paid_off_date" to debts table
		{
			ID: "202307191521",
			Migrate: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE debts ADD COLUMN forecast_paid_off_date VARCHAR(50);`,
				}

				return migration.ExecMultiple(tx, strings.Join(changes, " "))
			},
			Rollback: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE debts DROP COLUMN forecast_paid_off_date;`,
				}

				return migration.ExecMultiple(tx, strings.Join(changes, " "))
			},
		},
		// add column "forecast_start_investing_date", "forecast_financial_freedom_date", "forecast_millionaire_date" to sessions table
		// "forecast_emergency_budget_filled_date", "forecast_rainyday_budget_filled_date"
		{
			ID: "202307191630",
			Migrate: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE sessions ADD COLUMN forecast_emergency_budget_filled_date VARCHAR(50);`,
					`ALTER TABLE sessions ADD COLUMN forecast_rainyday_budget_filled_date VARCHAR(50);`,
					`ALTER TABLE sessions ADD COLUMN forecast_start_investing_date VARCHAR(50);`,
					`ALTER TABLE sessions ADD COLUMN forecast_financial_freedom_date VARCHAR(50);`,
					`ALTER TABLE sessions ADD COLUMN forecast_millionaire_date VARCHAR(50);`,
				}

				return migration.ExecMultiple(tx, strings.Join(changes, " "))
			},
			Rollback: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE sessions DROP COLUMN forecast_emergency_budget_filled_date,
						DROP COLUMN forecast_rainyday_budget_filled_date,
						DROP COLUMN forecast_start_investing_date,
						DROP COLUMN forecast_financial_freedom_date,
						DROP COLUMN forecast_millionaire_date;`,
				}

				return migration.ExecMultiple(tx, strings.Join(changes, " "))
			},
		},
		// replace name fun_fund to actual_fun_fund, add column expected_fun_fund to sessions table
		{
			ID: "202307221532",
			Migrate: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE sessions ADD COLUMN expected_fun_fund DOUBLE PRECISION;`,
					`ALTER TABLE sessions RENAME COLUMN fun_fund TO actual_fun_fund;`,
				}

				return migration.ExecMultiple(tx, strings.Join(changes, " "))
			},
			Rollback: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE sessions DROP COLUMN expected_fun_fund;`,
					`ALTER TABLE sessions RENAME COLUMN actual_fun_fund TO fun_fund;`,
				}

				return migration.ExecMultiple(tx, strings.Join(changes, " "))
			},
		},
		// add "forecast_bankrupt" to sessions table
		{
			ID: "202307251001",
			Migrate: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE sessions ADD COLUMN forecast_bankrupt VARCHAR(50);`,
				}

				return migration.ExecMultiple(tx, strings.Join(changes, " "))
			},
			Rollback: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE sessions DROP COLUMN forecast_bankrupt;`,
				}

				return migration.ExecMultiple(tx, strings.Join(changes, " "))
			},
		},
	})

	return nil
}
