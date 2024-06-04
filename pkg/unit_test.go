package pkg

import (
	"damir/internal/data"
	"damir/internal/entity"
	_ "fmt"
	_"net/http"
	_"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	_ "golang.org/x/crypto/bcrypt"
)

func TestInsertGame(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %s", err)
	}
	defer db.Close()
	app := &Application{
		Models: data.NewModels(db),
    }
	game := &entity.Game{
		Name: "Warcraft",
		Price: 1000,
		Genres: []string{"Adventure", "Horror"},
	}
	mock.ExpectQuery(`INSERT INTO games`).
		WithArgs(game.Name, game.Price, pq.Array(game.Genres)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "version"}).AddRow(1, time.Now(), 1))

	err = app.Models.Games.Insert(game)
	if err != nil {
		t.Errorf("unexpected error inserting games: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetGame(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %s", err)
	}
	defer db.Close()

	app := &Application{
		Models: data.NewModels(db),
	}

	mock.ExpectQuery(`SELECT \* FROM games WHERE`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "name", "price", "genres", "version"}).
			AddRow(1, time.Now(), "Warcraft", 1000, pq.Array([]string{"Adventure", "Horror"}), 1))

	game, err := app.Models.Games.Get(1)
	if err != nil {
		t.Errorf("unexpected error retrieving game: %s", err)
	}
	if game.ID != 1 {
		t.Errorf("incorrect game ID, got %d want %d", game.ID, 1)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteGame(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %s", err)
	}
	defer db.Close()

	app := &Application{
		Models: data.NewModels(db),
	}

	mock.ExpectExec(`DELETE FROM games WHERE`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = app.Models.Games.Delete(1)
	if err != nil {
		t.Errorf("unexpected error deleting game: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}


