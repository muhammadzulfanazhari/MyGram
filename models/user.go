package models

import (
	"MyGram/helpers"
	"errors"
	"strings"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type User struct {
	GormModel
	Username string  `gorm:"not null;uniqueIndex" json:"username" form:"username" valid:"required~Username is required"`
	Email    string  `gorm:"not null;uniqueIndex" json:"email" form:"email" valid:"required~email is required,email~Invalid email format"`
	Password string  `gorm:"not null" json:"password" form:"password" valid:"required~password is required,minstringlength(6)~Password has to have a minimum legth of characters"`
	Age      int     `gorm:"not null" json:"age" form:"age" valid:"required~Age is required,"`
	Photos   []Photo `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"photos"`
	// Comment  []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"comments"`
}

type ResUser struct {
	ID       uint   `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	errorMsg := ""
	_, errCreate := govalidator.ValidateStruct(u)

	if errCreate != nil {
		err = errCreate
		return
	}

	if u.Age <= 8 {
		errorMsg += "Age Must be greater than 8, "
	}

	if errorMsg != "" {
		err = errors.New(strings.TrimSuffix(errorMsg, ", "))
	}

	u.Password = helpers.HashPass(u.Password)
	// err = nil
	return
}
