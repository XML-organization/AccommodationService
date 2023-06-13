package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HostGrade struct {
	ID          uuid.UUID `json:"id"`
	UserId      uuid.UUID `json:"user_id"`
	UserName    string    `json:"user_name" `
	UserSurname string    `json:"user_surname"`
	Grade       float64   `json:"grade"`
}

func (hostGrade *HostGrade) BeforeCreate(scope *gorm.DB) error {
	hostGrade.ID = uuid.New()
	return nil
}
