package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidateRequestInput(s interface{}) []error {
	validate := validator.New()
	err := validate.Struct(s)

	if err != nil {
		errs := make([]error, 0)
		for _, e := range err.(validator.ValidationErrors) {
			errs = append(errs, fmt.Errorf("invalid field: %s", e.Field()))
		}
		return errs
	}

	return nil
}
