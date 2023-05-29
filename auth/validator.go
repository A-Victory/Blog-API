package auth

import (
	"github.com/A-Victory/blog-API/models"
	"github.com/go-playground/validator/v10"
)

type Validation struct {
	validate *validator.Validate
}

func NewValidator() *Validation {
	validate := validator.New()
	return &Validation{validate}
}

func (va *Validation) ValidateUserInfo(u models.User) error {
	err := va.validate.Struct(u)
	if err != nil {
		return err
	}
	return nil
}
