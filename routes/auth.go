package routes

import (
	"akmmp241/go-jwt/controllers"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterAuthRouters(r *mux.Router, c *controllers.AuthController) {
	router := r.PathPrefix("/auth").Subrouter()
	router.HandleFunc("/register", c.Register).Methods(http.MethodPost)
	router.HandleFunc("/login", c.Login).Methods(http.MethodPost)
}
