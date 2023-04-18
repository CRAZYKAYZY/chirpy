package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CRAZYKAYZY/chirpy/internal/config"
	"github.com/CRAZYKAYZY/chirpy/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
)

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Email           string `json:"email"`
	CurrentPassword string `json:"currentpassword"`
	NewPassword     string `json:"newpassword"`
}

type UserRes struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

func CreateUserHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req database.User

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := db.CreateUser(req.Email, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res := UserRes{
			ID:       user.ID,
			Email:    user.Email,
			Password: user.Password,
		}

		// Write the response
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetUsersHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		user, err := db.GetUser(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res := UserRes{
			ID:    user.ID,
			Email: user.Email,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetAllUsersHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := db.GetUsers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var res []UserRes
		for _, user := range users {
			res = append(res, UserRes{
				ID:    user.ID,
				Email: user.Email,
			})
		}

		respondWithJSON(w, http.StatusOK, res)
	}
}

func UpdateUsersHandler(db *database.DB, cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the claims
		claims := r.Context().Value("claims").(jwt.MapClaims)

		// Extract the user ID from the token claim's Subject
		userID, err := strconv.Atoi(claims["sub"].(string))
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Parse the request body
		var req UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Update the user in the database
		updatedUser, err := db.UpdateUser(userID, req.Email, req.CurrentPassword, req.NewPassword)
		if err != nil {
			fmt.Printf("error updating user: %v", err)
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		// Write the response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		res := UserRes{
			ID:    updatedUser.ID,
			Email: updatedUser.Email,
		}

		json.NewEncoder(w).Encode(res)
	}
}
