package pkg

import (
	"damir/internal/entity"
	"damir/internal/filters"
	rabbitmq "damir/internal/sender"
	"damir/internal/validator"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (app *Application) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Name     string `json:"name"`
        Surname  string `json:"surname"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    // Parse the request body into the anonymous struct.
    err := app.readJSON(w, r, &input)
    if err != nil {
        app.badRequestResponse(w, r, err)
        return
    }
    // Create a new user instance
    user := &entity.User{
        Name:      input.Name,
        Email:     input.Email,
        Surname:   input.Surname,
        Activated: true,
        Role:      "user",
    }
    // Use the Password.Set() method to generate and store the hashed password.
    err = user.Password.Set(input.Password)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }
    // Validate the user struct
    v := validator.New()
    if entity.ValidateUser(v, user); !v.Valid() {
        app.failedValidationResponse(w, r, v.Errors)
        return
    }
    // Insert the user data into the database
    err = app.Models.Users.Insert(user)
    if err != nil {
        switch {
        case errors.Is(err, entity.ErrDuplicateEmail):
            v.AddError("email", "a user with this email address already exists")
            app.failedValidationResponse(w, r, v.Errors)
        default:
            app.serverErrorResponse(w, r, err)
        }
        return
    }
    // Generate an activation token
    token, err := app.Models.Tokens.New(user.ID, 3*24*time.Hour, entity.ScopeActivation)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }
    // Send activation email
    app.background(func() {
        data := map[string]interface{}{
            "activationToken": token.Plaintext,
            "userID":          user.ID,
        }
        err := app.Mailer.Send(user.Email, "user_welcome.tmpl", data)
        if err != nil {
            app.Logger.PrintError(err, nil)
        }
    })
    message := fmt.Sprintf("User registered: %s", user.Email)
    err = rabbitmq.PublishMessage(message)
    if err != nil {
        app.Logger.PrintError(err, nil)
    }
    err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }
}

func (app *Application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate the plaintext token provided by the client.
	v := validator.New()
	if entity.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Retrieve the details of the user associated with the token using the
	// GetForToken() method (which we will create in a minute). If no matching record
	// is found, then we let the client know that the token they provided is not valid.
	user, err := app.Models.Users.GetForToken(entity.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Update the user's activation status.
	user.Activated = true
	// Save the updated user record in our database, checking for any edit conflicts in
	// the same way that we did for our movie records.
	err = app.Models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// If everything went successfully, then we delete all activation tokens for the
	// user.
	err = app.Models.Tokens.DeleteAllForUser(entity.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send the updated user details to the client in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) DeleteUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	fmt.Printf("the id %d", id)
	err = app.Models.Users.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "user successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	message := fmt.Sprintf("User deleted: %d", id)
	rabbitmq.PublishMessage(message)
}

func (app *Application) GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	user, err := app.Models.Users.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	message := fmt.Sprintf("User info retrived: %s", user.Name)
	rabbitmq.PublishMessage(message)
}

func (app *Application) GetAllUserInfoHandler(w http.ResponseWriter, r *http.Request) {
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
	users, err := app.Models.Users.GetAll(input.Name, input.Filters)
	//users, err := app.Models.Users.GetAll(input.Name, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"users": users}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	fmt.Fprintf(w, "%+v\n", input)
	message := fmt.Sprintf("All users info retrived: %d", len(users))
	rabbitmq.PublishMessage(message)
}

func (app *Application) EditUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	user, err := app.Models.Users.Get(id)
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
		Name      *string `json:"name"`
		Surname   *string `json:"surname"`
		Email     *string `json:"email"`
		Activated *bool   `json:"activated"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Surname != nil {
		user.Surname = *input.Surname
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Activated != nil {
		user.Activated = *input.Activated
	}

	v := validator.New()
	if entity.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.Models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	message := fmt.Sprintf("User updated: %d", user.ID)
	rabbitmq.PublishMessage(message)
}
