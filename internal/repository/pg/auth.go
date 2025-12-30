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

func (a *AuthRepository) GetUserByLogPass(user entity.DBUser) (userId uuid.UUID, err error) {
	if err := a.DB.Find(&userId, user.Login, user.Password).Error; err != nil {
		return uuid.Nil, err
	}
	return
}

func (a *AuthRepository) CreateUser(user *entity.User) error {
	if err := a.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}
