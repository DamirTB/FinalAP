package tests

import (
	"os"
	_"bytes"
	"damir/internal/data"
	"damir/internal/entity"
	"damir/pkg"
	_"encoding/json"
	"fmt"
	_"reflect"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/julienschmidt/httprouter"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
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

func TestDeleteUser(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("error creating mock database: %s", err)
    }
    defer db.Close()
    app := &pkg.Application{
        Models: data.NewModels(db),
    }
    userID := int64(1)
    mock.ExpectExec(`DELETE FROM user_info where id = \$1`).
        WithArgs(userID).
        WillReturnResult(sqlmock.NewResult(1, 1)) 
    err = app.Models.Users.Delete(userID)
    if err != nil {
        t.Errorf("unexpected error deleting user: %s", err)
    }
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectations: %s", err)
    }
}

func TestSetup(t *testing.T) {
	originalDSN := os.Getenv("DSN")
	originalEmail := os.Getenv("email")
	originalPassword := os.Getenv("password")

	os.Setenv("DSN", "test_dsn")
	os.Setenv("email", "test_email")
	os.Setenv("password", "test_password")

	defer func() {
		os.Setenv("DSN", originalDSN)
		os.Setenv("email", originalEmail)
		os.Setenv("password", originalPassword)
	}()

	cfg := pkg.Setup()

	assert.Equal(t, 4000, cfg.Port)
	assert.Equal(t, "development", cfg.Env)
	assert.Equal(t, "test_dsn", cfg.Db.Dsn)
	assert.Equal(t, 25, cfg.Db.MaxOpenConns)
	assert.Equal(t, 25, cfg.Db.MaxIdleConns)
	assert.Equal(t, "15m", cfg.Db.MaxIdleTime)
	assert.Equal(t, 2.0, cfg.Limiter.Rps)
	assert.Equal(t, 4, cfg.Limiter.Burst)
	assert.Equal(t, true, cfg.Limiter.Enabled)
	assert.Equal(t, "smtp.office365.com", cfg.Smtp.Host)
	assert.Equal(t, 587, cfg.Smtp.Port)
	assert.Equal(t, "test_email", cfg.Smtp.Username)
	assert.Equal(t, "test_password", cfg.Smtp.Password)
	assert.Equal(t, "Test <221363@astanait.edu.kz>", cfg.Smtp.Sender)
}

func TestInsertGameFailed(t *testing.T) {
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
    // Change the expectation to not return any rows
    mock.ExpectQuery(`INSERT INTO games`).
        WithArgs(game.Name, game.Price, pq.Array(game.Genres)).
        WillReturnRows(sqlmock.NewRows([]string{}))

    err = app.Models.Games.Insert(game)
    // Check if the error is not nil, as the insertion should fail
    if err == nil {
        t.Errorf("expected error while inserting game, but got nil")
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectations: %s", err)
    }
}
