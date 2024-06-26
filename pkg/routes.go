package pkg

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) Routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/games", app.GetAllGamesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/games", app.requireAdminUser(app.CreateGameHandler))
	router.HandlerFunc(http.MethodGet, "/v1/games/:id", app.ShowGameHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/games/:id", app.requireAdminUser(app.DeleteGameHandler))
	router.HandlerFunc(http.MethodPut, "/v1/games/:id", app.requireAdminUser(app.UpdateGameHandler))

	router.HandlerFunc(http.MethodPost, "/v1/order", app.requireActivatedUser(app.createOrderHandler))
	router.HandlerFunc(http.MethodGet, "/v1/orders", app.requireActivatedUser(app.getAllOrdersHandler))
	router.HandlerFunc(http.MethodPut, "/v1/orders", app.requireActivatedUser(app.refundOrderHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.RegisterUserHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users", app.GetAllUserInfoHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/users/:id", app.requireAdminUser(app.DeleteUserInfoHandler))
	router.HandlerFunc(http.MethodGet, "/v1/users/:id", app.requireActivatedUser(app.GetUserInfoHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/users/:id", app.requireAdminUser(app.EditUserInfoHandler))

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.CreateAuthenticationTokenHandler)
	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
