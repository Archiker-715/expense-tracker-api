package entity

import "github.com/google/uuid"

type Users struct {
	UserId   uuid.UUID `gorm:"unique"`
	Login    string    `gorm:"unique"`
	Password []byte
}

type UserAuthRegistration struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type DBUser struct {
	Login    string
	Password []byte
}
