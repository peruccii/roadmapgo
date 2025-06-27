package utils

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

func ValidateFields(input interface{}) error {
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return errors.New("invalid input" + err.Error())
	}
	return nil
}
