package data

import (
	r "damir/internal/usecase"
	"damir/internal/usecase/repo"
	"database/sql"
	"errors"
  )

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesч]n't exist in our database.
var (
	ErrRecordNotFound = errors.New("record (row, entry) not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Movies r.MovieRepository
	Users  r.UserRepository
	Tokens r.TokenRepository 
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies: repo.MovieModel{DB: db},
		Users:  repo.UserModel{DB: db},
		Tokens: repo.TokenModel{DB: db}, 
	}
}
