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

func (ord OrderModel) GetAll(user_id int32) ([]entity.Order, error) {
    var orders []entity.Order

    query := `
        SELECT id, game_id, order_date, status, version
        FROM orders
        WHERE user_id = $1`

    rows, err := ord.DB.Query(query, user_id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var order entity.Order
        err := rows.Scan(&order.ID, &order.GameID, &order.OrderDate, &order.Status, &order.Version)
        if err != nil {
            return nil, err
        }
        orders = append(orders, order)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return orders, nil
}
