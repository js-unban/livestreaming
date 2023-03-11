//go:build wireinject
// +build wireinject

package models

import (
	"github.com/google/wire"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"os"
)

var UserAlreadyExists = errors.New("user already exists")

type User struct {
	gorm.Model
	ID       uint `gorm:"primaryKey"`
	Name     string
	Email    string `gorm:"unique"`
	Password string
}

type UsersModel struct {
	db *gorm.DB
}

var SIGNING_KEY string

var Users UsersModel = initializeUsersModel()

func init() {
	log.SetReportCaller(true)

	ensureSigningKey()
}

func ensureSigningKey() {
	SIGNING_KEY = os.Getenv("SIGNING_KEY")
	if SIGNING_KEY == "" {
		log.Fatalln("SIGNING_KEY env variable can't be empty")
	}
}

func NewUsersDB(d *gorm.DB) UsersModel {
	d.AutoMigrate(&User{})
	return UsersModel{d}
}

func initializeUsersModel() UsersModel {

	wire.Build(NewUsersDB, NewGorm)
	return UsersModel{}
}

func (userModel *UsersModel) FindUserById(id uint) (*User, bool) {
	var user User
	if queryResult := userModel.db.First(&user, "id = ?", id); queryResult.Error != nil {
		if !errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			log.Error("db error found", queryResult.Error)
		}
		return nil, false
	} else {
		return &user, true
	}
}
func (userModel *UsersModel) FindUserByEmail(email string) (*User, bool) {
	var user User
	if queryResult := userModel.db.First(&user, "email = ?", email); queryResult.Error != nil {
		if !errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			log.Error("db error found", queryResult.Error)
		}
		return nil, false
	} else {
		return &user, true
	}
}

func hashPassword(password string) (string, error) {
	if newPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err != nil {
		//log.error
		return "", errors.New("Error generating password")
	} else {
		return string(newPassword), nil
	}
}

func (userModel *UsersModel) CreateUser(user *User) (*User, error) {
	if _, ok := userModel.FindUserByEmail(user.Email); ok {
		return nil, UserAlreadyExists
	}
	newPassword, err := hashPassword(user.Password)
	if err != nil {
		log.Error("user coudln't be created", err)
		return nil, errors.Wrap(err, "User couldn't be created")
	}
	user.Password = newPassword
	if queryResult := userModel.db.Create(user); queryResult.Error != nil {
		log.Error("user coudln't be created")

		return nil, errors.Wrap(queryResult.Error, "User couldn't be created")
	} else {
		return user, nil
	}
}

func (userModel *UsersModel) CreateOrGetUser(user *User) (*User, error) {
	newPassword, err := hashPassword(user.Password)
	if err != nil {
		// log.error
		return nil, errors.Wrap(err, "User couldn't be created")
	}
	user.Password = newPassword
	if queryResult := userModel.db.FirstOrCreate(user, "email = ?", user.Email); queryResult.Error != nil {
		return nil, errors.Wrap(queryResult.Error, "User couldn't be created")
	} else {
		return user, nil
	}
}

func (userModel *UsersModel) ValidateUserPassword(email string, password string) bool {
	if user, ok := userModel.FindUserByEmail(email); !ok {
		// log.debug
		return false
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return false
		}
		return true
	}
}
