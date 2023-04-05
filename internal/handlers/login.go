package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/CRAZYKAYZY/chirpy/internal/config"
	"github.com/CRAZYKAYZY/chirpy/internal/database"
	"github.com/golang-jwt/jwt/v5"
)

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func Login_handler(db *database.DB, cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginReq

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

		// Set expiration time for access token
		accessTokenExpirationTime := time.Now().Add(1 * time.Hour)

		//set jwt claims
		claims := jwt.MapClaims{
			"iss": "chirpy",
			"sub": strconv.Itoa(user.ID),
			"iat": jwt.NewNumericDate(time.Now().UTC()),
			"exp": jwt.NewNumericDate(accessTokenExpirationTime.UTC()),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString([]byte(cfg.JwtSecret))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		res := loginResponse{
			ID:    user.ID,
			Email: user.Email,
			Token: signedToken,
		}

		json.NewEncoder(w).Encode(res)
	}
}
