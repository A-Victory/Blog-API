package auth

import (
	"fmt"

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
		for _, err := range err.(validator.ValidationErrors) {
			return fmt.Errorf("please input %s, field is required!", err.Field())
		}
	}
	return nil
}
