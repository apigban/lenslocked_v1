package models

import (
	"errors"

	"github.com/apigban/lenslocked_v1/hash"
	"github.com/apigban/lenslocked_v1/rand"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrNotFound is returned when a resource is not found in the database
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an invalid ID is provided to a method like Delete()
	ErrInvalidID = errors.New("models: ID provided was invalid")

	ErrInvalidPassword = errors.New("models: incorrect password provided")
)

const userPwPepper = "peppa"
const hmacSecretKey = "secret-hmac-key"

// User represents the user model in the database
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"` //not going to be stored in the database
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"` //not going to be stored in the database
	RememberHash string `gorm:"not null;unique_index"`
}

// UserDB is used to interact with the users database.
type UserDB interface {
	// Methods for querying for single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// Methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// Used to close a DB connection
	Close() error

	// Migration helpers
	AutoMigrate() error
	DestructiveReset() error
}

func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := &userValidator{
		hmac:   hmac,
		UserDB: ug,
	}
	return &userService{
		UserDB: uv,
	}, nil
}

var _ UserDB = &userGorm{}

// UserService is a set of methods used to manipulate
// and work with the user model
type UserService interface {
	// Authenticate will pverify the provided email and password are correct.
	// If correct, the user associated to that email will be returned.
	// Can also return error:
	// ErrNotFound, ErrInvalidPassword, or catchall error
	Authenticate(email, password string) (*User, error)
	UserDB
}

// Implementation of the userService
type userService struct {
	UserDB
}

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

type userValidator struct {
	UserDB
	hmac hash.HMAC
}

// ByRemember will hash the remember token and then call
// ByRemember on the subsequent UserDB layer.
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}

	if err := runUserValFuncs(&user,
		uv.hmacRemember); err != nil {
		return nil, err
	}

	return uv.UserDB.ByRemember(user.RememberHash)
}

// Create will create the provided user and backfill the data
// like ID, CreatedAt and UpdatedAt
func (uv *userValidator) Create(user *User) error {
	// Always set remember has on Create(),
	if user.Remember != "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}

	err := runUserValFuncs(user, uv.bcryptPassword, uv.hmacRemember)
	if err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

// Update will hash a remember hash if token is provided
// in the user object
func (uv *userValidator) Update(user *User) error {
	err := runUserValFuncs(user, uv.bcryptPassword, uv.hmacRemember)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

// Delete will delete the user with the provided ID
// in the provided user object
func (uv *userValidator) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	return uv.UserDB.Delete(id)
}

// bcryptPassword will hash a user's password with a predefined pepper
// and bcrypt if the password field is not an empty string
func (uv *userValidator) bcryptPassword(user *User) error {
	// If password provided is empty, then user is not trying to set password
	if user.Password == "" {
		return nil
	}
	//TODO - Create validation function for password entry: tooShort, noUpper, noNumber, noSymbol
	pwBytes := []byte(user.Password + userPwPepper) // add pepper to password and cast concatenated string to byteslice
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = "" // Clear out password from memory, avoids logging to stdout
	return nil
}

// bcryptPassword will hash a user's password with a predefined pepper
// and bcrypt if the password field is not an empty string
func (uv *userValidator) hmacRemember(user *User) error {
	// If remember token provided is empty, no need to hash the token
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true) // TODO - remove when env == production
	return &userGorm{
		db: db,
	}, nil
}

type userGorm struct {
	db *gorm.DB
}

// ByID will look up a user by ID provided
// Case 1 - user, nil
// Case 2 - nil, ErrNotFound
// Case 3 - nil, otherError
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail will look up a user by Email Address provided
// Case 1 - user, nil
// Case 2 - nil, ErrNotFound
// Case 3 - nil, otherError
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// ByRemember finds a user by their remember token
// This method expects the remember token to be hashed
// Errors are the same as ByEmail
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Authenticate can be used to authenticate the user with the given user and password.
func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper))
	if err != nil { // if error IS nil, fallthrough
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}
	return foundUser, nil
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
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
	//TODO - create specific errors, like if it is invalid, or user already exists
}

// Delete will delete the user with the provided ID
// in the provided user object
func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// Update will update the provided user with all of the data
// in the provided user object
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

//Close closes the UserService data connection
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

//DestructiveReset drops the user table and rebuilds it
func (ug *userGorm) DestructiveReset() error {
	if err := ug.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return ug.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate the users table
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}
