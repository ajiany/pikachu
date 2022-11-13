package validators

import (
	"github.com/go-playground/validator/v10"
	"strings"

	"golang.org/x/net/publicsuffix"
)

var blackList []string

func init() {
	blackList = []string{
		"google.com",
		"baidu.com",
		"qq.com",
		"163.com",
		"facebook.com",
		"instagram.com",
	}
}

func domain(fl validator.FieldLevel) bool {
	d := fl.Field().String()

	d = strings.TrimSpace(strings.ToLower(d))

	registrableDomain, err := publicsuffix.EffectiveTLDPlusOne(d)
	if err != nil {
		return false
	}

	for _, bd := range blackList {
		if bd == registrableDomain {
			return false
		}
	}

	return true
}
