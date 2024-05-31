package config

import (
	"damir/internal/data"
	"damir/internal/mailer"
	"sync"
	"damir/internal/jsonlog"
	"flag"
	"os"
)

type application struct {
	config Config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

const version = "1.0.0"

func Setup() Config{
	var cfg Config

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.Db.Dsn, "db-dsn", os.Getenv("DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.Db.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.Db.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.Db.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max idle time")
	// flag.StringVar(&cfg.db.maxLifetime, "db-max-lifetime", "1h", "PostgreSQL max idle time")

	flag.Float64Var(&cfg.Limiter.Rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.Limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.Limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.Smtp.Host, "smtp-host", "smtp.office365.com", "SMTP host")
	flag.IntVar(&cfg.Smtp.Port, "smtp-port", 587, "SMTP port")
	flag.StringVar(&cfg.Smtp.Username, "smtp-username", os.Getenv("email"), "SMTP username")
	flag.StringVar(&cfg.Smtp.Password, "smtp-password", os.Getenv("password"), "SMTP password")
	flag.StringVar(&cfg.Smtp.Sender, "smtp-sender", "Test <221363@astanait.edu.kz>", "SMTP sender")
	flag.Parse()
	return cfg
}
