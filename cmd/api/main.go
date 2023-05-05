package main

import (
	"dullahan/config"

	dbutil "dullahan/internal/util/db"

	"github.com/M15t/ghoul/pkg/server"
)

func main() {
	cfg, err := config.Load()
	checkErr(err)

	db, err := dbutil.New(cfg.DbDsn, cfg.DbLog)
	checkErr(err)
	// connection.Close() is not available for GORM 1.20.0
	// defer db.Close()

	sqlDB, err := db.DB()
	defer sqlDB.Close()

	// Initialize HTTP server
	e := server.New(&server.Config{
		Stage:        cfg.Stage,
		Port:         cfg.Port,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		AllowOrigins: cfg.AllowOrigins,
		Debug:        cfg.Debug,
	})

	// Static page for Swagger API specs
	e.Static("/swaggerui", "swaggerui")

	// Start the HTTP server
	server.Start(e, cfg.Stage == "development")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
