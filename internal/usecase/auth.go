package usecase

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"student-cert-service/internal/domain"
)

type AuthRepository interface {
	CreateUser(username, password string) (int, error)
	GetUserByUsername(username string) (*domain.User, error)
}

type AuthUseCase struct {
	repo      AuthRepository
	jwtSecret string
}

func NewAuthUseCase(repo AuthRepository, jwtSecret string) *AuthUseCase {
	return &AuthUseCase{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (a *AuthUseCase) Register(username, password string) error {
	existingUser, err := a.repo.GetUserByUsername(username)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = a.repo.CreateUser(username, string(hashedPassword))
	return err
}

func (a *AuthUseCase) Login(username, password string) (string, error) {
	user, err := a.repo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(a.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
