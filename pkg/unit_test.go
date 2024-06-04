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
		t.Errorf("unexpected error inserting movie: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// func TestGetGame(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("error creating mock database: %s", err)
// 	}
// 	defer db.Close()
// 	app := &Application{
// 		Models: data.NewModels(db),
//     }
// 	mock.ExpectQuery(`SELECT * FROM movies WHERE id = $1`).
// 		WithArgs(1).
// 		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "name", "price", "genres", "version"}).
// 			AddRow(1, time.Now(), "Warcraft", 1000, pq.Array([]string{"Adventure", "Horror"}), 1))

// 	movie, err := app.Models.Games.Get(1)
// 	if err != nil {
// 		t.Errorf("unexpected error fetching movie: %s", err)
// 	}

// 	if movie.ID != 1 {
// 		t.Errorf("incorrect movie ID, got %d, want %d", movie.ID, 1)
// 	}

// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 	}
// }