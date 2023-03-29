package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

// create a request struct type
type chirpReqValidate struct {
	Body string `json:"body"`
}

// create a respond struct type
type ResChirpValidate struct {
	Valid     bool   `json:"valid"`
	CleanBody string `json:"cleaned_body,omitempty"`
	Error     string `json:"error,omitempty"`
}

// handler controller that takes in the instance of the request and then performs the logic
func ValidateChirpHandler(w http.ResponseWriter, r *http.Request) {
	//request instance
	var req chirpReqValidate

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.Body) > 140 {
		validateRes := ResChirpValidate{
			Valid: false,
			Error: "Chirp is too long",
		}
		json.NewEncoder(w).Encode(validateRes)
		return
	}

	validateRes := ResChirpValidate{
		Valid:     true,
		CleanBody: Profane(req.Body),
	}

	json.NewEncoder(w).Encode(validateRes)

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
