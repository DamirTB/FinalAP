package pkg

import (
	"damir/internal/data"
	"damir/internal/jsonlog"
	"damir/internal/mailer"
	cfn "damir/pkg/config"
	"flag"
	"os"
	"sync"
)

type Applicaiton struct {
	Config cfn.Config
	Logger *jsonlog.Logger
	Models data.Models
	Mailer mailer.Mailer
	Wg     sync.WaitGroup
}

const version = "1.0.0"

func Setup() cfn.Config {
	var cfg cfn.Config

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.Db.Dsn, "db-dsn", os.Getenv("DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.Db.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.Db.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.Db.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max idle time")

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

