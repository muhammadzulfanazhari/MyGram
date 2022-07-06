package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Photo struct {
	GormModel
	Title     string `gorm:"not null" json:"title" form:"title" valid:"required~Title of your photo is required"`
	Caption   string `json:"caption"`
	Photo_url string `gorm:"not null" json:"photo_url" form:"photo_url" valid:"required~URL of your photo is required"`
	UserID    uint   `json:"user_id"`
	User      *User
}

type ResPhoto struct {
	ID        uint       `json:"id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitemtpy"`
	Title     string     `json:"title,omitempty"`
	Caption   string     `json:"caption,omitempty"`
	Photo_url string     `json:"photo_url,omitempty"`
	UserID    uint       `json:"user_id,omitempty"`
	User      *ResUser   `json:"user,omitempty"`
}

func (p *Photo) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(p)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}
