package pkg

import (
  "damir/internal/entity"
  "damir/internal/filters"
  rabbitmq "damir/internal/sender"
  "damir/internal/validator"
  "errors"
  "fmt"
  "net/http"
  "strings"
)

func (app *Application) CreateGameHandler(w http.ResponseWriter, r *http.Request) {
  var input struct {
    Name   string   `json:"name"`
    Price  int32    `json:"price"`
    Genres []string `json:"genres"`
  }
  err := app.readJSON(w, r, &input)
  if err != nil {
    app.errorResponse(w, r, http.StatusBadRequest, err.Error())
  }

  game := &entity.Game{
    Name:   input.Name,
    Price:  input.Price,
    Genres: input.Genres,
  }

  err = app.Models.Games.Insert(game)
  if err != nil {
    app.serverErrorResponse(w, r, err)
    return
  }

  err = app.writeJSON(w, http.StatusCreated, envelope{"game": game}, nil)
  if err != nil {
    app.serverErrorResponse(w, r, err)
  }
  message := fmt.Sprintf("Game created: GameID=%d, GameName=%s", game.ID, game.Name)
  rabbitmq.PublishMessage(message)
  // fmt.Fprintf(w, "%+v\n", input) //+v here is adding the field name of a value // https://pkg.go.dev/fmt
}

func (app *Application) ShowGameHandler(w http.ResponseWriter, r *http.Request) {
  id, err := app.readIDParam(r)
  if err != nil {
    app.notFoundResponse(w, r)
  }

  game, err := app.Models.Games.Get(id)
  if err != nil {
    switch {
    case errors.Is(err, entity.ErrRecordNotFound):
      app.notFoundResponse(w, r)
    default:
      app.serverErrorResponse(w, r, err)
    }
    return
  }
  err = app.writeJSON(w, http.StatusOK, envelope{"game": game}, nil)
  if err != nil {
    app.serverErrorResponse(w, r, err)
  }
  message := fmt.Sprintf("Game retrived: %d", game.ID)
  rabbitmq.PublishMessage(message)
}

func (app *Application) DeleteGameHandler(w http.ResponseWriter, r *http.Request) {
  id, err := app.readIDParam(r)
  if err != nil {
    app.notFoundResponse(w, r)
    return
  }
  err = app.Models.Games.Delete(id)
  if err != nil {
    switch {
    case errors.Is(err, entity.ErrRecordNotFound):
      app.notFoundResponse(w, r)
    default:
      app.serverErrorResponse(w, r, err)
    }
    return
  }
  err = app.writeJSON(w, http.StatusOK, envelope{"message": "game successfully deleted"}, nil)
  if err != nil {
    app.serverErrorResponse(w, r, err)
  }
  message := fmt.Sprintf("Game deleted: %d", id)
  rabbitmq.PublishMessage(message)
}

func (app *Application) UpdateGameHandler(w http.ResponseWriter, r *http.Request) {
  id, err := app.readIDParam(r)
  if err != nil {
    app.notFoundResponse(w, r)
    return
  }
  game, err := app.Models.Games.Get(id)
  if err != nil {
    switch {
    case errors.Is(err, entity.ErrRecordNotFound):
      app.notFoundResponse(w, r)
    default:
      app.serverErrorResponse(w, r, err)
    }
    return
  }

  var input struct {
    Name   string   `json:"name"`
    Price  int32    `json:"price"`
    Genres []string `json:"genres"`
  }

  err = app.readJSON(w, r, &input)
  if err != nil {
    app.serverErrorResponse(w, r, err)
    return
  }
  game.Name = input.Name
  game.Price = input.Price
  game.Genres = input.Genres
  err = app.Models.Games.Update(game)
  if err != nil {
    app.serverErrorResponse(w, r, err)
    return
  }
  err = app.writeJSON(w, http.StatusOK, envelope{"game": game}, nil)
  if err != nil {
    app.serverErrorResponse(w, r, err)
  }
  message := fmt.Sprintf("Game updated: %d", id)
  rabbitmq.PublishMessage(message)
}



func (app *Application) GetAllGamesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
	  Name string
	  filters.Filters
	}
	v := validator.New()
	qs := r.URL.Query()
	input.Name = app.readString(qs, "name", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "-id", "-name"}
	if filters.ValidateFilters(v, input.Filters); !v.Valid() {
	  app.failedValidationResponse(w, r, v.Errors)
	  return
	}
	games, err := app.Models.Games.GetAll(input.Name, input.Filters)
	//users, err := app.Models.Users.GetAll(input.Name, input.Filters)
	if err != nil {
	  app.serverErrorResponse(w, r, err)
	  return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"games": games}, nil)
	if err != nil {
	  app.serverErrorResponse(w, r, err)
	}
	fmt.Fprintf(w, "%+v\n", input)
	// Format the games into a readable string
	var gameDetails []string
	for _, game := range games {
	  gameDetails = append(gameDetails, fmt.Sprintf("{ID: %d, Name: %s, Price: %d, Genres: %v}", game.ID, game.Name, game.Price, game.Genres))
	}
	logMessage := fmt.Sprintf("All games retrieved: [%s]", strings.Join(gameDetails, ", "))
  
	// Log the formatted message
	rabbitmq.PublishMessage(logMessage)
  }