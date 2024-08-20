package controllers

import (
	"akmmp241/go-jwt/configs"
	"akmmp241/go-jwt/helpers"
	"akmmp241/go-jwt/models/domains"
	"akmmp241/go-jwt/models/dto"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type AuthController struct {
	DB        *sql.DB
	Validator *validator.Validate
	Config    *configs.Config
}

func NewAuthController(DB *sql.DB, validator *validator.Validate, config *configs.Config) *AuthController {
	return &AuthController{DB: DB, Validator: validator, Config: config}
}

func (c AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var register dto.Register
	if err := json.NewDecoder(r.Body).Decode(&register); err != nil {
		helpers.Response(w, http.StatusInternalServerError, "Invalid request body", nil)
		return
	}
	defer r.Body.Close()

	err := c.Validator.Struct(register)
	if err != nil && errors.As(err, &validator.ValidationErrors{}) {
		helpers.Response(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	hashedPassword, err := helpers.HashPassword(register.Password)
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, "Failed to hash password", nil)
		return
	}

	user := domains.User{
		Name:     register.Name,
		Email:    register.Email,
		Password: hashedPassword,
	}

	SQL := "INSERT INTO users (id, name, email, password) VALUES (NULL, ?, ?, ?)"
	_, err = c.DB.ExecContext(r.Context(), SQL, user.Name, user.Email, user.Password)
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, "Failed to save users to database", nil)
		return
	}

	helpers.Response(w, http.StatusCreated, "User saved successfully", nil)
	return
}

func (c AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var login dto.Login
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		helpers.Response(w, http.StatusInternalServerError, "Invalid request body", nil)
		return
	}
	defer r.Body.Close()

	err := c.Validator.Struct(login)
	if err != nil && errors.As(err, &validator.ValidationErrors{}) {
		helpers.Response(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var user domains.User
	SQL := "SELECT id, name, email, password FROM users WHERE email = ?"
	rows, err := c.DB.QueryContext(r.Context(), SQL, login.Email)
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, "Failed to save users to database", nil)
		return
	}
	defer rows.Close()
	if !rows.Next() {
		helpers.Response(w, http.StatusUnauthorized, "Wrong Credentials", nil)
		return
	}
	err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if isVerified := helpers.VerifyPassword(login.Password, user.Password); !isVerified {
		helpers.Response(w, http.StatusUnauthorized, "Wrong Credentials", nil)
		return
	}

	key := c.Config.C.GetString("JWT_KEY")
	token, err := helpers.CreateToken(&user, []byte(key))
	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, "Failed to create token", nil)
		return
	}

	helpers.Response(w, http.StatusCreated, "Successfully Login", token)
}
