package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"sync"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	if err := db.ensureDB(); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) ensureDB() error {
	_, err := os.Stat(db.path)
	if os.IsNotExist(err) {
		fmt.Println("Db does not exist, creating new...")
		return db.writeDB(DBStructure{Chirps: make(map[int]Chirp)})
	}
	return err
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	jsonData, err := json.Marshal(dbStructure)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(db.path, jsonData, 0644)
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := ioutil.ReadFile(db.path)

	if err != nil {
		return DBStructure{}, err
	}

	dbStructure := DBStructure{}

	if err := json.Unmarshal(data, &dbStructure); err != nil {
		return DBStructure{}, err
	}

	return dbStructure, nil
}

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
