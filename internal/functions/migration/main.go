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
					NetAssets      float64 `json:"net_assets"`

					ActualEmergencyFund   float64 `json:"actual_emergency_fund"`
					ExpectedEmergencyFund float64 `json:"expected_emergency_fund"`
					ActualRainydayFund    float64 `json:"actual_rainyday_fund"`
					ExpectedRainydayFund  float64 `json:"expected_rainyday_fund"`
					FunFund               float64 `json:"fun_fund"`
					Investment            float64 `json:"investment"`
					RetirementPlan        float64 `json:"retirement_plan"`

					IsAchivedEmergencyFund  bool `json:"is_achived_emergency_fund"`
					IsAchivedRainydayFund   bool `json:"is_achived_rainyday_fund"`
					IsAchivedInvestment     bool `json:"is_achived_investment"`
					IsAchivedRetirementPlan bool `json:"is_achived_retirement_plan"`

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

					DebtPaidOffEachMonth     float64 `json:"debt_paid_off_each_month"`
					InterestPaidOffEachMonth float64 `json:"interest_paid_off_each_month"`
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
				return tx.Exec(`ALTER TABLE sessions ADD COLUMN total_all_expense DOUBLE DEFAULT 0 AFTER total_all_income;`).Error
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
					`ALTER TABLE sessions ADD COLUMN total_monthly_payment_debt DOUBLE DEFAULT 0 AFTER total_all_expense;`,
					`ALTER TABLE sessions DROP COLUMN monthly_payment_debt;`,
				}

				return migration.ExecMultiple(tx, strings.Join(changes, " "))
			},
			Rollback: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE sessions ADD COLUMN monthly_payment_debt DOUBLE DEFAULT 0;`,
					`ALTER TABLE sessions DROP COLUMN total_monthly_payment_debt;`,
				}

				return migration.ExecMultiple(tx, strings.Join(changes, " "))
			},
		},
	})

	return nil
}
