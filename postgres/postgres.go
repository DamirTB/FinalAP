package postgres

import (
	"database/sql"
	"time"
	"context"
)

func OpenDB(cfg string, maxIdleCon int, maxOpenCon int, maxIdletime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(maxIdleCon)
	db.SetMaxOpenConns(maxOpenCon)

	duration, err := time.ParseDuration(maxIdletime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	
	// lifetime, err := time.ParseDuration(cfg.db.maxIdleTime)
	// if err != nil {
	// 	return nil, err
	// }
	// db.SetConnMaxLifetime(lifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx) 

	if err != nil {
		return nil, err
	}

	return db, nil
}