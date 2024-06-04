package tests

import (
    "database/sql"
    "testing"
	"damir/internal/entity"
	"damir/internal/data"
	"damir/pkg"
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
        t.Fatalf("Failed to get movie after update: %s", err)
    }
    if updatedGame.Name != "World of warcraft Updated" {
        t.Errorf("Updated game name %s does not match expected title 'World of warcraft Updated'", updatedGame.Name)
    }
    err = app.Models.Games.Delete(game.ID)
    if err != nil {
        t.Fatalf("Failed to delete game: %s", err)
    }
}
