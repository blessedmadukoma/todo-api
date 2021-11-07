package handlers

import (
	"errors"
	"fmt"
	"go-gorm-pg/models"
	"log"
	"os"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// HashPassword hashes password from user input
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // 14 is the cost for hashing the password.
	return string(bytes), err
}

// CheckPasswordHash checks password hash and password from user input if they match
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("password incorrect")
	}
	return nil
}

// Prepare: strips user input of any white spaces
func Prepare(u *models.User) {
	u.Email = strings.TrimSpace(u.Email)
	u.Firstname = strings.TrimSpace(u.Firstname)
	u.Lastname = strings.TrimSpace(u.Lastname)
}

// Validate user input
func Validate(u *models.User, action string) error {
	switch strings.ToLower(action) {
	case "login":
		if u.Email == "" {
			return errors.New("email is required")
		}
		if u.Password == "" {
			return errors.New("password is required")
		}
		return nil
	default: // this is for creating a user, where all fields are required
		if u.Firstname == "" {
			return errors.New("first name is required")
		}
		if u.Lastname == "" {
			return errors.New("last name is required")
		}
		if u.Email == "" {
			return errors.New("email is required")
		}
		if u.Password == "" {
			return errors.New("password is required")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("invalid Email")
		}
		return nil
	}
}

// GetUser returns a user based on email
func GetUser(u *models.User, db *gorm.DB) (*models.User, error) {
	account := &models.User{}
	if err := db.Debug().Table("users").Where("email = ?", u.Email).First(account).Error; err != nil {
		return nil, err
	}
	fmt.Println(account)
	return account, nil
}

// GetAllUsers returns a list of all the user
func GetAllUsers(db *gorm.DB) (*[]models.User, error) {
	users := []models.User{}
	if err := db.Debug().Table("users").Find(&users).Error; err != nil {
		return &[]models.User{}, err
	}
	return &users, nil
}

func EncodeAuthToken(uid uint) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Could not load .env in encode auth")
	}
	claims := jwt.MapClaims{}
	claims["userID"] = uid
	claims["IssuedAt"] = time.Now().Unix()
	claims["ExpiresAt"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	return token.SignedString([]byte(os.Getenv("SECRET")))
}
