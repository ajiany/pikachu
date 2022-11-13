package form

import (
	"fmt"

	"github.com/ajiany/pikachu/validators"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"gopkg.in/go-playground/validator.v9"
)

var serverLangs = []language.Tag{
	language.AmericanEnglish,
	language.Chinese,
}
var matcher = language.NewMatcher(serverLangs)

type Form interface {
	Do(c *gin.Context) (interface{}, error)
}

func DoThisForm(form Form, c *gin.Context) (bool, interface{}) {
	headerLocale := c.GetHeader("X-LOCALE")
	locale, _ := language.MatchStrings(matcher, headerLocale)

	if err := c.ShouldBind(form); err != nil {
		errs, ok := err.(validator.ValidationErrors)

		if ok {
			msg := []string{}
			for _, e := range errs {
				msg = append(msg, e.Translate(validators.GetTrans(locale.String())))
			}
			c.JSON(422, gin.H{"errors": msg})
			return false, nil
		} else {
			fmt.Printf("400 request : %v\n", err)
			c.String(400, "")
			return false, nil
		}
	}

	ret, err := form.Do(c)
	if err == nil {
		return true, ret
	}

	var (
		code              int
		message           string
		statuscode        int
		isWithStatusError bool
		wse               *WithStatusError
		isCustomError     bool
		ce                *CustomError
		isWithCodeError   bool
		wce               *WithCodeError
	)

	if IsUnprocessableError(err) {
		code = 422
	} else if IsUnauthorizedError(err) {
		code = 401
	} else if IsInternalStateError(err) {
		code = 520
	} else if IsNotFoundError(err) {
		code = 404
	} else if isCustomError, ce = IsCustomError(err); isCustomError {
		code = ce.Code
		statuscode = ce.Status
	} else if isWithCodeError, wce = IsWithCodeError(err); isWithCodeError {
		code = wce.Code
	} else if isWithStatusError, wse = IsWithStatusError(err); isWithStatusError {
		code = 200
		statuscode = wse.Status
	} else if IsContextCancel(err) {
		code = 408
	} else {
		panic(err)
	}

	message = err.Error()

	if isCustomError || isWithStatusError {
		c.JSON(code, gin.H{
			"status": statuscode,
			"msg":    message,
		})
	} else {
		c.JSON(code, gin.H{
			"errors": []string{message},
		})
	}
	return false, nil

}
