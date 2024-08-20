package main

import (
	"akmmp241/go-jwt/configs"
	"akmmp241/go-jwt/controllers"
	"akmmp241/go-jwt/middlewares"
	"akmmp241/go-jwt/routes"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	config := configs.NewConfig()

	db := configs.ConnectDB(config)
	v := validator.New()

	authController := controllers.NewAuthController(db, v, config)
	userController := controllers.NewUserController(db, v, config)

	authMiddleware := middlewares.NewAuthMiddleware(config)

	r := mux.NewRouter()
	router := r.PathPrefix("/api").Subrouter()
	routes.RegisterAuthRouters(router, authController)
	routes.RegisterUserRoutes(router, authMiddleware, userController)

	log.Fatal(http.ListenAndServe(`:3000`, router))
}
