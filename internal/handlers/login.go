package handlers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/CRAZYKAYZY/chirpy/internal/database"
)

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func Login_handler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req database.User

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := db.GetUserByEmail(req.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			http.Error(w, "invalid password", http.StatusUnauthorized)
			return
		}

		// Write the response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		res := loginResponse{
			ID:    user.ID,
			Email: user.Email,
		}

		json.NewEncoder(w).Encode(res)
	}
}
