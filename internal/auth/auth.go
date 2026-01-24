package auth

import (
	"errors"
	"fmt"
	"os"

	"github.com/Archiker-715/expense-tracker-api/internal/entity"
	"github.com/Archiker-715/expense-tracker-api/internal/repository/pg"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	repo *pg.AuthRepository
}

func NewAuthService(repo *pg.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

var (
	secretKey = []byte(os.Getenv("JWT_KEY"))
)

func (a *AuthService) Authorization(user entity.UserAuthRegistration) (entity.AuthResponse, error) {
	DBUser := entity.DBUser{
		Login:    user.Login,
		Password: Encode(user.Password),
	}
	userId, err := a.repo.GetUserByLogPass(DBUser)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.AuthResponse{}, errors.New("User does not exist or password is incorrect. Try again")
		}
		return entity.AuthResponse{}, fmt.Errorf("get user: %w", err)
	}

	out, err := a.generateToken(userId)
	if err != nil {
		return entity.AuthResponse{}, fmt.Errorf("get token: %w", err)
	}

	return out, nil
}

func (a *AuthService) Registration(user entity.UserAuthRegistration) error {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("create UUID: %w", err)
	}
	newUser := entity.Users{
		UserId:   newUUID,
		Login:    user.Login,
		Password: Encode(user.Password),
	}

	err = a.repo.CreateUser(&newUser)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("User %q already exists", user.Login)
		}
		return fmt.Errorf("DB error: %w", err)
	}

	return nil
}
