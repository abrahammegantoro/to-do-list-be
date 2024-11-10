package user

import (
	"context"
	"os"
	"time"

	"github.com/abrahammegantoro/to-do-list-be/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Login(ctx context.Context, username string, password string) (domain.User, error)
	Register(ctx context.Context, user domain.User) error
	GetByUsername(ctx context.Context, email string) (res domain.User, err error)
	GetByID(ctx context.Context, id int64) (res domain.User, err error)
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(ur UserRepository) *UserService {
	return &UserService{
		userRepository: ur,
	}
}

func (u *UserService) Login(ctx context.Context, auth *domain.AuthCredentials) (res domain.User, token string, err error) {
	user, err := u.userRepository.GetByUsername(ctx, auth.Username)
	if err != nil {
		return domain.User{}, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(auth.Password))
	if err != nil {
		return domain.User{}, "", domain.ErrCredential
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	token, err = claims.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return domain.User{}, "", domain.ErrInternalServerError
	}

	return user, token, nil
}

func (u *UserService) Register(ctx context.Context, user *domain.User) (token string, err error) {
	_, err = u.userRepository.GetByUsername(ctx, user.Username)
	if err == nil {
		return "", domain.ErrUsernameTaken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", domain.ErrInternalServerError
	}

	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err = u.userRepository.Register(ctx, *user)
	if err != nil {
		return "", domain.ErrInternalServerError
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	token, err = claims.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", domain.ErrInternalServerError
	}

	return token, nil
}
