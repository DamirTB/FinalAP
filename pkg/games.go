package pkg

import (
	"damir/internal/entity"
	"errors"
	_ "fmt"
	"net/http"
)

func (app *Application) createGameHandler(w http.ResponseWriter, r *http.Request) {
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
	// fmt.Fprintf(w, "%+v\n", input) //+v here is adding the field name of a value // https://pkg.go.dev/fmt
}

func (app *Application) showGameHandler(w http.ResponseWriter, r *http.Request) {
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
}

func (app *Application) deleteGameHandler(w http.ResponseWriter, r *http.Request) {
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
}

// func(app *Application) updateGameHandler(w http.ResponseWriter, r *http.Request){
// 	id, err := app.readIDParam(r)
// 	if err != nil{ 
// 		app.notFoundResponse(w, r)
// 		return
// 	}
// 	game, err := app.Models.Games.Get(id)
// }


func (app *Application) updateGameHandler(w http.ResponseWriter, r *http.Request) {
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
		Name   		string   `json:"name"`
		Price    	int32    `json:"price"`
		Genres  	[]string `json:"genres"`
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
}
