package pg

import (
	"github.com/Archiker-715/expense-tracker-api/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (a *AuthRepository) GetUserByLogPass(user entity.DBUser) (uuid.UUID, error) {
	var dbUser entity.Users
	if err := a.DB.Where("login = ? AND password = ?", user.Login, user.Password).First(&dbUser).Error; err != nil {
		return uuid.Nil, err
	}
	return dbUser.UserId, nil
}

func (a *AuthRepository) CreateUser(user *entity.Users) error {
	if err := a.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}
