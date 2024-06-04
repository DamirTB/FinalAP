package pkg

import (
  "damir/internal/entity"
  rabbitmq "damir/internal/sender"
  "errors"
  "fmt"
  "net/http"
)

func (app *Application) createOrderHandler(w http.ResponseWriter, r *http.Request) {
  var input struct {
    Game_id int64 `json:"game_id"`
  }
  err := app.readJSON(w, r, &input)
  if err != nil {
    app.errorResponse(w, r, http.StatusBadRequest, err.Error())
  }
  game, err := app.Models.Games.Get(input.Game_id)
  if err != nil {
    switch {
    case errors.Is(err, entity.ErrRecordNotFound):
      app.notFoundResponse(w, r)
    default:
      app.serverErrorResponse(w, r, err)
    }
    return
  }
  user := app.contextGetUser(r)
  Order := &entity.Order{
    GameID: input.Game_id,
    UserID: user.ID,
  }
  if user.Balance < game.Price {
    app.failedPayment(w, r)
    return
  }
  err = app.Models.Users.PayBalance(game.Price, user)
  if err != nil {
    app.serverErrorResponse(w, r, err)
    return
  }
  err = app.Models.Orders.Insert(int32(Order.GameID), int32(Order.UserID), Order)
  if err != nil {
    app.serverErrorResponse(w, r, err)
    return
  }
  err = app.writeJSON(w, http.StatusCreated, envelope{"game": game, "order": Order}, nil)
  if err != nil {
    app.serverErrorResponse(w, r, err)
  }
  message := fmt.Sprintf("Order created: GameID=%d, UserID=%d", Order.GameID, Order.UserID)
  rabbitmq.PublishMessage(message)
}

func (app *Application) getAllOrdersHandler(w http.ResponseWriter, r *http.Request) {
    user := app.contextGetUser(r)
    if user == nil {
        app.errorResponse(w, r, http.StatusUnauthorized, "user not authenticated")
        return
    }

    orders, err := app.Models.Orders.GetAll(int32(user.ID))
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }

    err = app.writeJSON(w, http.StatusOK, orders, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }
}
