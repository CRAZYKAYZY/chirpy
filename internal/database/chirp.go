package database

import (
	"sort"
)

func (db *DB) CreateChirp(body string) (Chirp, error) {
	// Lock the database for writing
	db.mux.Lock()
	defer db.mux.Unlock()

	// Load the database into memory
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	// Generate a new ID for the chirp
	newID := len(dbStructure.Chirps) + 1

	// Create the new chirp
	newChirp := Chirp{
		ID:   newID,
		Body: body,
	}

	// Add the new chirp to the database
	dbStructure.Chirps[newID] = newChirp

	// Save the updated database to disk
	if err := db.writeDB(dbStructure); err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	// Lock the database for reading
	db.mux.RLock()
	defer db.mux.RUnlock()

	// Load the database into memory
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	// Convert the map to a slice and sort by ID
	var chirps []Chirp
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	// Sort the chirps by ID
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	return chirps, nil
}
