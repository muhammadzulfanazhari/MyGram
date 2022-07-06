package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Comment struct {
	GormModel
	Message string `gorm:"not null" json:"message" form:"message" valid:"required~message is required"`
	UserID  uint   `json:"user_id"`
	User    *User  `json:"User,omitempty"`
	PhotoID uint   `json:"photo_id"`
	Photo   *Photo `json:"Photo,omitempty"`
}

type ResComment struct {
	ID        uint       `json:"id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitemtpy"`
	Message   string     `json:"message,omitempty"`
	UserID    uint       `json:"user_id,omitempty"`
	User      *ResUser   `json:"User,omitempty"`
	PhotoID   uint       `json:"photo_id,omitempty"`
	Photo     *ResPhoto  `json:"Photo,omitempty"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(c)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}
