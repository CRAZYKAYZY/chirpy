package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/CRAZYKAYZY/chirpy/internal/database"
)

// create a respond struct type
type ChirpRes struct {
	ID    int    `json:"id"`
	Body  string `json:"body"`
	Valid bool   `json:"valid"`
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
			ID:   chirp.ID,
			Body: Profane(chirp.Body),
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
