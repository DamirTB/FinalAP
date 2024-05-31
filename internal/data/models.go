package data

import (
	r "damir/internal/usecase"
	"damir/internal/usecase/repo"
	"database/sql"
	"errors"
  )

var (
	ErrRecordNotFound = errors.New("record (row, entry) not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Users  r.UserRepository
	Tokens r.TokenRepository 
	Games  r.GameRepository
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:  repo.UserModel{DB: db},
		Tokens: repo.TokenModel{DB: db}, 
		Games: repo.GameModel{DB: db},
	}
}
