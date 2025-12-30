package entity

import "github.com/google/uuid"

type User struct {
	UserId   uuid.UUID `gorm:"unique"`
	Login    string    `gorm:"unique"`
	Password [32]byte
}

type UserAuthRegistration struct {
	Login    string
	Password string
}

type DBUser struct {
	Login    string
	Password [32]byte
}
