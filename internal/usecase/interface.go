package usecase

import (
	"damir/internal/entity"
    "damir/internal/filters"
	"time"
)

type UserRepository interface {
	Insert(user *entity.User) error
	Get(id int64) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	GetForToken(tokenScope, tokenPlaintext string) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id int64) error
	GetAll(name string, filters filters.Filters) ([]*entity.User, error)
	PayBalance(price int32, user *entity.User) error
}

type TokenRepository interface{
	Insert(token *entity.Token) error
	New(userID int64, ttl time.Duration, scope string) (*entity.Token, error)
	DeleteAllForUser(scope string, userID int64) error
}

type GameRepository interface {
    Insert(game *entity.Game) error
    Get(id int64) (*entity.Game, error)
    Delete(id int64) error
	Update(game *entity.Game) error
	GetAll(name string, filters filters.Filters) ([]*entity.Game, error)
}

type OrderRepository interface{
	Insert(game_id int32, user_id int32, order *entity.Order) error
}