package validators

import (
	"gopkg.in/go-playground/validator.v9"
)

func password(fl validator.FieldLevel) bool {
	pw := fl.Field().String()

	if len(pw) < 6 || len(pw) > 20 {
		return false
	}

	//TODO
	return true
}
