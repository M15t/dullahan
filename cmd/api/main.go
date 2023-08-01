package main

import (
	"context"
	"dullahan/config"
	"time"

	"dullahan/internal/api/v1/auth"
	"dullahan/internal/api/v1/customer/debt"
	"dullahan/internal/api/v1/customer/expense"
	"dullahan/internal/api/v1/customer/income"
	"dullahan/internal/api/v1/customer/session"
	"dullahan/internal/db"
	"dullahan/internal/rbac"
	"dullahan/internal/util/crypter"
	dbutil "dullahan/internal/util/db"

	"github.com/M15t/ghoul/pkg/server"
	"github.com/M15t/ghoul/pkg/server/middleware/jwt"
	"github.com/allegro/bigcache/v3"

	_ "dullahan/internal/util/swagger" // Swagger stuffs
)

func main() {
	cfg, err := config.Load()
	checkErr(err)

	gdb, err := dbutil.New(cfg.DbDsn, cfg.DbLog)
	checkErr(err)
	// connection.Close() is not available for GORM 1.20.0
	// defer db.Close()

	sqlDB, err := gdb.DB()
	defer sqlDB.Close()

	// * Initialize HTTP server
	e := server.New(&server.Config{
		Stage:        cfg.Stage,
		Port:         cfg.Port,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		AllowOrigins: cfg.AllowOrigins,
		Debug:        cfg.Debug,
	})

	// * Static page for Swagger API specs
	e.Static("/swaggerui", "swaggerui")

	// * Initialize services
	dbSvc := db.New(gdb)
	rbacSvc := rbac.New(cfg.Debug)
	crypterSvc := crypter.New()
	jwtSvc := jwt.New(cfg.JwtAlgorithm, cfg.JwtSecret, cfg.JwtDuration)

	cache, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(1*time.Minute))
	defer cache.Reset()

	authSvc := auth.New(dbSvc, jwtSvc, crypterSvc, cfg)

	incomeSvc := income.New(dbSvc, rbacSvc, crypterSvc)
	expenseSvc := expense.New(dbSvc, rbacSvc, crypterSvc)
	debtSvc := debt.New(dbSvc, rbacSvc, crypterSvc)
	sessionSvc := session.New(dbSvc, rbacSvc, crypterSvc, cache)

	// * Initialize v1 API
	v1Router := e.Group("/v1")

	// * Initialize auth API
	auth.NewHTTP(authSvc, v1Router)

	// * Load jwt middleware
	v1cRouter := v1Router.Group("/customer")
	v1cRouter.Use(jwtSvc.MWFunc())

	income.NewHTTP(incomeSvc, authSvc, v1cRouter.Group("/incomes"))
	expense.NewHTTP(expenseSvc, authSvc, v1cRouter.Group("/expenses"))
	debt.NewHTTP(debtSvc, authSvc, v1cRouter.Group("/debts"))
	session.NewHTTP(sessionSvc, authSvc, v1cRouter.Group("/me"))

	// Start the HTTP server
	server.Start(e, cfg.Stage == "development")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
