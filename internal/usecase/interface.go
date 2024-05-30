package usecase

import (
	"damir/internal/entity"
    "damir/internal/filters"
	"time"
)

type MovieRepository interface {
    Insert(movie *entity.Movie) error
    Get(id int64) (*entity.Movie, error)
    Update(movie *entity.Movie) error
    Delete(id int64) error
}

type UserRepository interface {
	Insert(user *entity.User) error
	Get(id int64) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	GetForToken(tokenScope, tokenPlaintext string) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id int64) error
	GetAll(name string, filters filters.Filters) ([]*entity.User, error)
}

type TokenRepository interface{
	Insert(token *entity.Token) error
	New(userID int64, ttl time.Duration, scope string) (*entity.Token, error)
	DeleteAllForUser(scope string, userID int64) error
}
