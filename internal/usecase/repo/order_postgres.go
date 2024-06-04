package repo

import (
	"damir/internal/entity"
	"database/sql"
)

type OrderModel struct {
	DB *sql.DB
}

func (ord OrderModel) Insert(game_id int32, user_id int32, order *entity.Order) error {
	query := `
		INSERT INTO orders(user_id, game_id, status)
		VALUES ($1, $2, 'Accepted')
		RETURNING id, order_date, version`
	return ord.DB.QueryRow(query, game_id, user_id).Scan(&order.ID, &order.OrderDate, &order.Version)
}