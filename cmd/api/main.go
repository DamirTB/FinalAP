package main

import (
	"damir/internal/data"
	"damir/internal/jsonlog"
	"damir/internal/mailer"
	"damir/postgres"
	"os"
	"sync"
	cfn "damir/config"
)

type application struct {
	config cfn.Config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	cfg := cfn.Setup()
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := postgres.OpenDB(cfg.Db.Dsn, cfg.Db.MaxIdleConns, cfg.Db.MaxOpenConns, cfg.Db.MaxIdleTime)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.Sender),
	}
	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
