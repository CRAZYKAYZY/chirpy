package database

import (
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func (db *DB) CreateUser(email, password string) (User, error) {

	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	// Check if user with the same email already exists
	for _, user := range dbStructure.Users {
		if user.Email == email {
			return User{}, fmt.Errorf("user with email %s already exists", email)
		}
	}

	id := len(dbStructure.Users) + 1

	// Generate hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	// Create a new user with hashed password
	user := User{
		ID:       id,
		Email:    email,
		Password: string(hashedPassword),
	}

	dbStructure.Users[id] = user

	if err := db.writeDB(dbStructure); err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbstructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	for _, user := range dbstructure.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, errors.New("user not found")
}

func (db *DB) GetUser(id int) (User, error) {
	dbstructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	for _, user := range dbstructure.Users {
		if user.ID == id {
			return user, nil
		}
	}
	return User{}, errors.New("user is not found")
}

// GetUsers returns all users in the database
func (db *DB) GetUsers() ([]User, error) {

	// Load the database into memory
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	var users []User
	for _, user := range dbStructure.Users {
		users = append(users, user)
	}

	return users, nil
}

func (db *DB) UpdateUser(userID int, email, password string) (User, error) {
	user, err := db.GetUser(userID)
	if err != nil {
		return User{}, err
	}

	id := user.ID

	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	// Replace the user at that index with the updated user
	user.Email = email
	user.Password = string(hashedPassword)
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		log.Fatal(err)
	}

	return user, nil
}
