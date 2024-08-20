package routes

import (
	"akmmp241/go-jwt/controllers"
	"akmmp241/go-jwt/middlewares"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterUserRoutes(r *mux.Router, auth *middlewares.AuthMiddleware, c *controllers.UserController) {
	router := r.PathPrefix("/users").Subrouter()
	router.Use(auth.JWTAuth)
	router.HandleFunc("", c.Me).Methods(http.MethodGet)
}
