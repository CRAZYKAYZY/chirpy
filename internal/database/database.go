package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
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
		return db.writeDB(DBStructure{
			Chirps: make(map[int]Chirp),
			Users:  make(map[int]User),
		})
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
	db.mux.RLock()
	defer db.mux.RUnlock()

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
