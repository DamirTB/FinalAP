package pkg

import (
  "damir/internal/entity"
  rabbitmq "damir/internal/sender"
  "damir/internal/validator"
  "errors"
  "fmt"
  "net/http"
  "time"
)

func (app *Application) CreateAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
  var input struct {
    Email    string `json:"email"`
    Password string `json:"password"`
  }
  err := app.readJSON(w, r, &input)
  if err != nil {
    app.badRequestResponse(w, r, err)
    return
  }
  v := validator.New()
  entity.ValidateEmail(v, input.Email)
  entity.ValidatePasswordPlaintext(v, input.Password)
  if !v.Valid() {
    app.failedValidationResponse(w, r, v.Errors)
    return
  }
  user, err := app.Models.Users.GetByEmail(input.Email)
  if err != nil {
    switch {
    case errors.Is(err, entity.ErrRecordNotFound):
      app.invalidCredentialsResponse(w, r)
    default:
      app.serverErrorResponse(w, r, err)
    }
    return
  }
  match, err := user.Password.Matches(input.Password)
  if err != nil {
    app.serverErrorResponse(w, r, err)
    return
  }
  if !match {
    app.invalidCredentialsResponse(w, r)
    return
  }
  token, err := app.Models.Tokens.New(user.ID, 24*time.Hour, entity.ScopeAuthentication)
  if err != nil {
    app.serverErrorResponse(w, r, err)
    return
  }
  err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
  if err != nil {
    app.serverErrorResponse(w, r, err)
  }
  message := fmt.Sprintf("User authenticated: %s, User ID: %d", user.Email, user.ID)
  err = rabbitmq.PublishMessage(message)
  if err != nil {
    app.Logger.PrintError(err, nil)
  }
}