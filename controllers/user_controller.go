package controllers

import (
	"akmmp241/go-jwt/configs"
	"akmmp241/go-jwt/helpers"
	"akmmp241/go-jwt/models/dto"
	"database/sql"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type UserController struct {
	DB        *sql.DB
	Validator *validator.Validate
	Config    *configs.Config
}

func (c UserController) Me(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*helpers.MyCustomClaims)

	userDto := dto.UserDto{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}

	helpers.Response(w, http.StatusOK, "Success get data user", &userDto)
}

func NewUserController(DB *sql.DB, validator *validator.Validate, config *configs.Config) *UserController {
	return &UserController{DB: DB, Validator: validator, Config: config}
}
