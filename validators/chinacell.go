package validators

import (
	"regexp"

	"gopkg.in/go-playground/validator.v9"
)

const REG_CHINA_CELL = `^1[3456789]\d{9}$`

var reg *regexp.Regexp

func chinaCell(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}
	if reg == nil {
		var err error
		reg, err = regexp.Compile(REG_CHINA_CELL)
		if err != nil {
			panic(err)
		}
	}
	return reg.MatchString(fl.Field().String())
}
