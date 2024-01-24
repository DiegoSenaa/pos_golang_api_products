package handlers

import (
	"encoding/json"
	"github.com/DiegoSenaa/golang-api/internal/dto"
	"github.com/DiegoSenaa/golang-api/internal/entity"
	"github.com/DiegoSenaa/golang-api/internal/infra/database"
	"github.com/go-chi/jwtauth"
	"net/http"
	"time"
)

type UserHandler struct {
	UserDB database.UserInterface
}

func NewUserHandler(db database.UserInterface) *UserHandler {
	return &UserHandler{
		UserDB: db,
	}
}

type Error struct {
	Message string
}

// CreateUser godoc
// @Summary      Create user
// @Description  Create user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request     body      dto.CreateUserRequest  true  "user request"
// @Success      201
// @Failure      400         {object}  Error "internal Server Error"
// @Failure      500         {object}  Error "internal Server Error"
// @Router       /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errResponse := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errResponse)
		return
	}
	err = h.UserDB.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetUserByEmail godoc
// @Summary      Get user by email
// @Description  Get user by email
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        email   query     string  true   "email to search"  Format(email)
// @Success      200  {object}  entity.User
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /users [get]
func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.UserDB.FindByEmail(email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// GetJwt godoc
// @Summary      Get a user JWT
// @Description  Get a user JWT
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request   body     dto.UserTokenRequest  true  "user credentials"
// @Success      200  {object}  dto.UserTokenResponse
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /users/generate_token [post]
func (h *UserHandler) GetJwt(w http.ResponseWriter, r *http.Request) {
	var user dto.UserTokenRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errResponse := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	u, err := h.UserDB.FindByEmail(user.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		errResponse := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errResponse)
		return
	}
	if !u.ValidatePassword(user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpirationTime := r.Context().Value("jwtExpirationTime").(int)

	_, token, _ := jwt.Encode(map[string]interface{}{
		"sub": u.Id.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpirationTime)).Unix(),
	})

	response := dto.UserTokenResponse{AccessToken: token}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
