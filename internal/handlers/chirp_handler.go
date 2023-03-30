package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/CRAZYKAYZY/chirpy/internal/database"
	"github.com/go-chi/chi/v5"
)

// create a respond struct type
type ChirpRes struct {
	ID    int    `json:"id"`
	Body  string `json:"body"`
	Valid bool   `json:"valid,omitempty"`
	Error string `json:"error,omitempty"`
}

func CreateChirpHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req database.Chirp
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		chirp, err := db.CreateChirp(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(req.Body) > 140 {
			validateRes := ChirpRes{
				Valid: false,
				Error: "chirp too long",
			}
			json.NewEncoder(w).Encode(validateRes)
			return
		}

		res := ChirpRes{
			Valid: true,
			Body:  Profane(chirp.Body),
		}

		json.NewEncoder(w).Encode(res)

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetChirpsHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirps, err := db.GetChirps()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var res []ChirpRes
		for _, chirp := range chirps {
			res = append(res, ChirpRes{
				ID:   chirp.ID,
				Body: chirp.Body,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetChirpHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}
		chirps, err := db.GetChirps()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, chirps[id-1])
	}
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func Profane(Chirp string) string {
	profaneWords := []string{"shit", "punk", "ass", "drugs", "guns", "bITCHES"}

	cleanWords := make([]string, 0)
	words := strings.Split(Chirp, " ")

	for _, word := range words {
		lowercaseWord := strings.ToLower(strings.Trim(word, "!.,?"))

		for _, pfWrd := range profaneWords {
			lowerPfWrd := strings.ToLower(pfWrd)
			if lowercaseWord == lowerPfWrd {
				word = "****"
				break
			}
		}
		cleanWords = append(cleanWords, word)
	}
	return strings.Join(cleanWords, " ")
}
