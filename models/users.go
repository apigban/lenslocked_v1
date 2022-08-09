package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound is returned when a resource is not found in the database
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an invalid ID is provided to a method like Delete()
	ErrInvalidID = errors.New("models: ID provided was invalid")
)

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true) // TODO - remove when env == production
	return &UserService{
		db: db,
	}, nil
}

type UserService struct {
	db *gorm.DB
}

// ByID will look up a user by ID provided
// Case 1 - user, nil
// Case 2 - nil, ErrNotFound
// Case 3 - nil, otherError
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail will look up a user by Email Address provided
// Case 1 - user, nil
// Case 2 - nil, ErrNotFound
// Case 3 - nil, otherError
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// first will query the provided gorm.DB and it will
// get the first item returned and place it to dst. If
// nothing is found the query, it will return ErrNotFound
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// Create will create the provided user and backfill the data
// like ID, CreatedAt and UpdatedAt
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
	//TODO - create specific errors, like if it is invalid, or user already exists
}

// Delete will delete the user with the provided ID
// in the provided user object
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error

}

// Update will update the provided user with all of the data
// in the provided user object
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

//Close closes the UserService data connection
func (us *UserService) Close() error {
	return us.db.Close()
}

//DestructiveReset drops the user table and rebuilds it
func (us *UserService) DestructiveReset() {
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})
}

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}
