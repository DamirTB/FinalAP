package entity

import (
	"time"
)

type Order struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	GameID    int64     `json:"game_id"`
	OrderDate time.Time `json:"order_date"`
	Status    string    `json:"status"`
	Version   int32     `json:"version"`
}
