package main

import (
  "damir/internal/data"
  "damir/internal/jsonlog"
  "damir/internal/mailer"
  rabbitmq "damir/internal/sender"
  cfn "damir/pkg"
  "damir/postgres"
  "os"
)

func main() {
  cfg := cfn.Setup()
  logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

  db, err := postgres.OpenDB(cfg.Db.Dsn, cfg.Db.MaxIdleConns, cfg.Db.MaxOpenConns, cfg.Db.MaxIdleTime)
  if err != nil {
    logger.PrintFatal(err, nil)
  }
  defer db.Close()
  logger.PrintInfo("database connection pool established", nil)

  err = rabbitmq.InitRabbitMQ()
  if err != nil {
    logger.PrintFatal(err, nil)
  }
  defer rabbitmq.CloseRabbitMQ()
  logger.PrintInfo("RabbitMQ connection established", nil)

  app := &cfn.Application{
    Config: cfg,
    Logger: logger,
    Models: data.NewModels(db),
    Mailer: mailer.New(cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.Sender),
  }

  err = app.Serve()
  if err != nil {
    logger.PrintFatal(err, nil)
  }
}