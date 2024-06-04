package tests

import (
    "database/sql"
    "testing"
	"damir/internal/entity"
	"damir/internal/data"
	"damir/pkg"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
	"net/http/httptest"
	"fmt"
	"encoding/json"
	"damir/internal/mailer"
	"os"
	"damir/internal/jsonlog"
	"bytes"
)


func setupDB(t *testing.T) *sql.DB {
    dsn := "postgres://postgres:12345@localhost:5432/fake?sslmode=disable"
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        t.Fatal("Failed to open a DB connection: ", err)
    }
    return db
}

func TestUserModelMethods(t *testing.T) {
    db := setupDB(t)
    defer db.Close()

    _, err := db.Exec("DELETE FROM user_info")
    if err != nil {
        t.Fatalf("Failed to clear test database: %s", err)
    }
	
	app := &pkg.Application{
        Models: data.NewModels(db),
    }
    password := "securepassword123"
    user := &entity.User{
        Name:      "John",
        Surname:   "Doe",
        Email:     "john.doe@example.com",
        Activated: true,
    }

    err = user.Password.Set(password)
    if err != nil {
        t.Fatalf("Failed to set password: %s", err)
    }

    err = app.Models.Users.Insert(user)
    if err != nil {
        t.Errorf("Failed to insert user: %s", err)
    }

    retrievedUser, err := app.Models.Users.Get(user.ID)
    if err != nil {
        t.Errorf("Failed to get user: %s", err)
    } else if retrievedUser.Email != user.Email {
        t.Errorf("Get user email %s does not match expected email %s", retrievedUser.Email, user.Email)
    }

    user.Name = "Updated John"
    err = app.Models.Users.Update(user)
    if err != nil {
        t.Errorf("Failed to update user: %s", err)
    }

    updatedUser, err := app.Models.Users.Get(user.ID)
    if err != nil {
        t.Errorf("Failed to get user after update: %s", err)
    } else if updatedUser.Name != "Updated John" {
        t.Errorf("Update user name %s does not match expected name 'Updated John'", updatedUser.Name)
    }

    err = app.Models.Users.Delete(user.ID)
    if err != nil {
        t.Errorf("Failed to delete user: %s", err)
    }
}

func TestGameModelMethods(t *testing.T) {
    db := setupDB(t)
    defer db.Close()

    _, err := db.Exec("DELETE FROM games")
    if err != nil {
        t.Fatalf("Failed to clear test database: %s", err)
    }
    app := &pkg.Application{
        Models: data.NewModels(db),
    }
    game := &entity.Game{
        Name:   	"World of warcraft",
        Price:    	2000,
        Genres:  	[]string{"MMORPG"},
    }
    err = app.Models.Games.Insert(game)
    if err != nil {
        t.Fatalf("Failed to insert game: %s", err)
    }
    retrievedGame, err := app.Models.Games.Get(game.ID)
    if err != nil {
        t.Fatalf("Failed to get game: %s", err)
    }
    if retrievedGame.Name != game.Name {
        t.Errorf("Retrieved game name %s does not match expected name %s", retrievedGame.Name, game.Name)
    }
    game.Name = "World of warcraft Updated"
    err = app.Models.Games.Update(game)
    if err != nil {
        t.Fatalf("Failed to update game: %s", err)
    }
    updatedGame, err := app.Models.Games.Get(game.ID)
    if err != nil {
        t.Fatalf("Failed to get game after update: %s", err)
    }
    if updatedGame.Name != "World of warcraft Updated" {
        t.Errorf("Updated game name %s does not match expected title 'World of warcraft Updated'", updatedGame.Name)
    }
    err = app.Models.Games.Delete(game.ID)
    if err != nil {
        t.Fatalf("Failed to delete game: %s", err)
    }
}

func TestUserHandlers(t *testing.T) {
	db := setupDB(t)
	if db == nil {
		t.Fatal("Database connection is nil")
	}
	defer db.Close()

	// smtpConfig := smtp{
	// 	host: "smtp.office365.com",
	// 	port: 587, 
	// 	username: os.Getenv("email"),
	// 	password: os.Getenv("password"),
	// 	sender: "Test <221363@astanait.edu.kz>",
	// }

	// appConfig := config{
	// 	smtp: smtpConfig,
	// }
	app := &pkg.Application{
		Models: data.NewModels(db),
		Mailer: mailer.New("smtp.office365.com", 587, "221363@astanait.edu.kz", "HAdCVB6quqmU3", "Test <221363@astanait.edu.kz>"),
		Logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
	}

	_, err := db.Exec(`DELETE FROM user_info `)
	if err != nil {
		t.Fatalf("Failed to clear test database: %s", err)
	}

	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, "/v1/users", app.RegisterUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.CreateAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users/:id", app.GetUserInfoHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/users/:id", app.EditUserInfoHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/users/:id", app.DeleteUserInfoHandler)

	body := strings.NewReader(`{"name":"John","surname":"Doe","email":"john.doe@example.com","password":"password123"}`)
	req, _ := http.NewRequest(http.MethodPost, "/v1/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Errorf("Expected status %d; got %d", http.StatusAccepted, w.Code)
	}

	var response struct {
		User struct {
			ID int64 `json:"id"`
		} `json:"user"`
	}
	err = json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if response.User.ID == 0 {
		t.Errorf("Expected a valid user ID; got %d", response.User.ID)
	}
	body = strings.NewReader(`{"email":"john.doe@example.com","password":"password123"}`)
	req, _ = http.NewRequest(http.MethodPost, "/v1/tokens/authentication", body)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d; got %d", http.StatusCreated, w.Code)
	}
	var tokenResponse struct {
		AuthenticationToken struct {
			Plaintext string `json:"Plaintext"`
		} `json:"authentication_token"`
	}
	err = json.NewDecoder(w.Body).Decode(&tokenResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if tokenResponse.AuthenticationToken.Plaintext == "" {
		t.Errorf("Expected a valid token; got an empty string")
	}

	fmt.Printf("Extracted Token: %s\n", tokenResponse.AuthenticationToken.Plaintext)

	
	// bearer_key := fmt.Sprintf("Bearer %s", tokenResponse.AuthenticationToken.Plaintext)
	//log.Printf("Params received: %+v", bearer_key)
	userURL := fmt.Sprintf("/v1/users/%d", response.User.ID)
	req, _ = http.NewRequest(http.MethodGet, userURL, nil)
	//req.Header.Set("Authorization", bearer_key)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, w.Code)
	}
	updateBody := strings.NewReader(`{"name":"John","surname":"Doe","email":"john.doe@example.com","password":"password123"}`)
	req, _ = http.NewRequest(http.MethodPatch, userURL, updateBody)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, w.Code)
	}
	req, _ = http.NewRequest(http.MethodDelete, userURL, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, w.Code)
	}
}


func TestGameHandlers(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	_, err := db.Exec("DELETE FROM games")
	if err != nil {
		t.Fatalf("Failed to clear test database: %s", err)
	}

	app := &pkg.Application{
		Models: data.NewModels(db),
	}

	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/v1/games", app.CreateGameHandler)
	router.HandlerFunc(http.MethodGet, "/v1/games/:id", app.ShowGameHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/games/:id", app.UpdateGameHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/games/:id", app.DeleteGameHandler)
	router.HandlerFunc(http.MethodGet, "/v1/games", app.GetAllGamesHandler)
	createGamePayload := []byte(`{"name":"Test Game", "price": 100, "genres":["Action", "Adventure"]}`)
	createGameRequest, _ := http.NewRequest(http.MethodPost, "/v1/games", bytes.NewBuffer(createGamePayload))
	createGameRequest.Header.Set("Content-Type", "application/json")
	createGameResponse := httptest.NewRecorder()
	router.ServeHTTP(createGameResponse, createGameRequest)

	if createGameResponse.Code != http.StatusCreated {
		t.Errorf("Expected status %d; got %d", http.StatusCreated, createGameResponse.Code)
	}
	var createGameResponseData struct {
		Game entity.Game `json:"game"`
	}
	err = json.NewDecoder(createGameResponse.Body).Decode(&createGameResponseData)
	if err != nil {
		t.Fatalf("Failed to decode create game response: %v", err)
	}

	getGameRequest, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/games/%d", createGameResponseData.Game.ID), nil)
	getGameResponse := httptest.NewRecorder()
	router.ServeHTTP(getGameResponse, getGameRequest)

	if getGameResponse.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, getGameResponse.Code)
	}
	deleteGameRequest, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/games/%d", createGameResponseData.Game.ID), nil)
	deleteGameResponse := httptest.NewRecorder()
	router.ServeHTTP(deleteGameResponse, deleteGameRequest)

	if deleteGameResponse.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, deleteGameResponse.Code)
	}
	getAllGamesRequest, _ := http.NewRequest(http.MethodGet, "/v1/games", nil)
	getAllGamesResponse := httptest.NewRecorder()
	router.ServeHTTP(getAllGamesResponse, getAllGamesRequest)

	if getAllGamesResponse.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, getAllGamesResponse.Code)
	}
}

func TestGameHandlers_Failed(t *testing.T) {
	db := setupDB(t)
	if db == nil {
		t.Fatal("Database connection is nil")
	}
	defer db.Close()
	_, err := db.Exec("DELETE FROM games")
	if err != nil {
		t.Fatalf("Failed to clear test database: %s", err)
	}
	app := &pkg.Application{
		Models: data.NewModels(db),
	}
	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, "/v1/games", app.CreateGameHandler)
	invalidGamePayload := []byte(`{"name":"", "price":"word", "genres":["Action", "Adventure"]}`)
	req, _ := http.NewRequest(http.MethodPost, "/v1/games", bytes.NewBuffer(invalidGamePayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d; got %d", http.StatusBadRequest, w.Code)
	}

	var createErrorResponse struct {
		ErrorMessage string `json:"error"`
	}
	err = json.NewDecoder(w.Body).Decode(&createErrorResponse)
	if err != nil {
		t.Fatalf("Failed to decode create game response: %s", err)
	}

	expectedCreateErrorMessage := `body contains incorrect JSON type for field "price"`
	if createErrorResponse.ErrorMessage != expectedCreateErrorMessage {
		t.Errorf("Expected error message '%s'; got '%s'", expectedCreateErrorMessage, createErrorResponse.ErrorMessage)
	}
	req, _ = http.NewRequest(http.MethodGet, "/v1/games/999", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d; got %d", http.StatusNotFound, w.Code)
	}
}

