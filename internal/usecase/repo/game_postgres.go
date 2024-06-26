package repo

import (
	"damir/internal/entity"
	"database/sql"
	"github.com/lib/pq"
	"errors"
	"damir/internal/filters"
	"context"
	"time"
	"fmt"
)

type GameModel struct {
	DB *sql.DB
}

func (g GameModel) Insert(game *entity.Game) error {
	query := `
		INSERT INTO games(name, price, genres)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, version`

	return g.DB.QueryRow(query, &game.Name, &game.Price, pq.Array(&game.Genres)).Scan(&game.ID, &game.CreatedAt, &game.Version)
}

func (m GameModel) Get(id int64) (*entity.Game, error) {
	if id < 1 {
		return nil, entity.ErrRecordNotFound
	}
	query := `
		SELECT *
		FROM games
		WHERE id = $1`
	var game entity.Game
	err := m.DB.QueryRow(query, id).Scan(
		&game.ID,
		&game.CreatedAt,
		&game.Name,
		&game.Price,
		pq.Array(&game.Genres),
		&game.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, entity.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &game, nil
}


func (g GameModel) Delete(id int64) error {
	if id < 1 {
		return entity.ErrRecordNotFound
	}
	query := `
		DELETE FROM games
		WHERE id = $1`
	result, err := g.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return entity.ErrRecordNotFound
	}
	return nil
}

func (g GameModel) Update(game *entity.Game) error {
	query := `
		UPDATE games
		SET name = $1, price = $2, genres = $3, version = version + 1
		WHERE id = $4
		RETURNING version`

	args := []interface{}{
		game.Name,
		game.Price,
		pq.Array(game.Genres),
		game.ID,
	}

	return g.DB.QueryRow(query, args...).Scan(&game.Version)
}

func (g GameModel) GetAll(name string, filters filters.Filters) ([]*entity.Game, error) {
	query := fmt.Sprintf(`
	SELECT *
	FROM games
	WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
	ORDER BY %s %s, id ASC
	LIMIT $2 OFFSET $3`, filters.SortColumn(), filters.SortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []any{name, filters.Limit(), filters.Offset()}
	rows, err := g.DB.QueryContext(ctx, query, args...)	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []*entity.Game

	for rows.Next() {
		var game entity.Game
		err := rows.Scan(
			&game.ID,
			&game.CreatedAt,
			&game.Name,
			&game.Price,
			pq.Array(&game.Genres),
			&game.Version,
		)
		if err != nil {
			return nil, err
		}
		games = append(games, &game)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return games, nil
}