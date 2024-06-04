package tests

import (
	"damir/internal/data"
	"damir/internal/entity"
	"damir/pkg"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"strings"
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
	app := &pkg.Application{
		Models: data.NewModels(db),
	}
	game := &entity.Game{
		Name:   "Warcraft",
		Price:  1000,
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

	app := &pkg.Application{
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

	app := &pkg.Application{
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

func TestUpdateMovie(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %s", err)
	}
	defer db.Close()
	app := &pkg.Application{
		Models: data.NewModels(db),
	}
	game := &entity.Game{
		ID:     1,
		Name:   "Updated Warcraft",
		Price:  1500,
		Genres: []string{"Adventure", "Horror"},
	}

	mock.ExpectQuery(`UPDATE games SET`).
		WithArgs(game.Name, game.Price, pq.Array(game.Genres), game.ID).
		WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow(2))

	err = app.Models.Games.Update(game)
	if err != nil {
		t.Errorf("unexpected error updating movie: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %s", err)
	}
	defer db.Close()
	app := &pkg.Application{
		Models: data.NewModels(db),
	}
	password := "password123"

	user := &entity.User{
		Name:      "John",
		Surname:   "Doe",
		Email:     "john.doe@example.com",
		Activated: true,
	}
	err = user.Password.Set(password)
	if err != nil {
		t.Fatalf("error hashing password: %s", err)
	}
	mock.ExpectQuery(`INSERT INTO user_info`).
		WithArgs(user.Name, user.Surname, user.Email, user.Password.Hash, user.Activated).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "version"}).AddRow(1, time.Now(), 1))

	err = app.Models.Users.Insert(user)
	if err != nil {
		t.Errorf("unexpected error inserting user: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %s", err)
	}
	defer db.Close()
	app := &pkg.Application{
		Models: data.NewModels(db),
	}

	user := &entity.User{
		ID:        1,
		Name:      "Jane",
		Surname:   "Doe",
		Email:     "jane.doe@example.com",
		Activated: true,
	}

	mock.ExpectQuery(`UPDATE user_info SET`).
		WithArgs(user.Name, user.Surname, user.Email, user.Activated, user.ID).
		WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow(2))

	err = app.Models.Users.Update(user)
	if err != nil {
		t.Errorf("unexpected error updating user: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllUsersHandler(t *testing.T) {
	app := &pkg.Application{}

	req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)

	rr := httptest.NewRecorder()

	app.Routes().ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d; got %d", http.StatusUnauthorized, rr.Code)
	}
}


func TestGetUserInfoHandler_Unathorized(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("error creating mock database: %s", err)
    }
    defer db.Close()

    app := &pkg.Application{
        Models: data.NewModels(db),
    }
    user := &entity.User{
        ID:        1,
        Name:      "Jane",
        Surname:   "Doe",
        Email:     "jane.doe@example.com",
        Activated: true,
    }
    mock.ExpectQuery(`SELECT \* FROM user_info WHERE`).
        WithArgs(user.ID).
        WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "email", "activated"}).
            AddRow(user.ID, user.Name, user.Surname, user.Email, user.Activated))

    req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/users/%d", user.ID), nil)
    rr := httptest.NewRecorder()
	app.Routes().ServeHTTP(rr, req)
    if rr.Code != http.StatusUnauthorized {
        t.Errorf("Expected status %d; got %d", http.StatusUnauthorized, rr.Code)
    }
    expectedBody := `{"error":"you must be authenticated to access this resource"}`
	fmt.Print(rr.Body.String())
    if !strings.Contains(rr.Body.String(), expectedBody) {
        t.Errorf("Expected body to contain %q", expectedBody)
    }
}


