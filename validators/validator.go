package validators

import (
	"reflect"
	"sync"

	"github.com/gin-gonic/gin/binding"
	enlocale "github.com/go-playground/locales/en_US"
	zhlocale "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/go-playground/validator.v9/translations/en"
	"gopkg.in/go-playground/validator.v9/translations/zh"
)

type DefaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var (
	ZhTrans ut.Translator
	EnTrans ut.Translator
)

func init() {
	zhTrans := zhlocale.New()
	enTrans := enlocale.New()
	uniTrans := ut.New(enTrans, zhTrans, enTrans)
	ZhTrans, _ = uniTrans.GetTranslator("zh")
	EnTrans, _ = uniTrans.GetTranslator("en_US")
}

var _ binding.StructValidator = &DefaultValidator{}

func (v *DefaultValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {

		v.lazyinit()

		if err := v.validate.Struct(obj); err != nil {
			return error(err)
		}
	}

	return nil
}

func GetTrans(locale string) ut.Translator {
	if locale == "zh" {
		return ZhTrans
	} else if locale == "en-US" {
		return EnTrans
	}
	return ZhTrans
}

func (v *DefaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *DefaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")
		zh.RegisterDefaultTranslations(v.validate, ZhTrans)
		en.RegisterDefaultTranslations(v.validate, EnTrans)

		// add any custom validations etc. here
		v.validate.RegisterValidation("chinacell", chinaCell)
		v.validate.RegisterValidation("password", password)

		v.validate.RegisterTranslation("required", ZhTrans, func(ut ut.Translator) error {
			return ut.Add("required", "{0}必填", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required", fe.Field())
			return t
		})
		v.validate.RegisterTranslation("chinacell", ZhTrans, func(ut ut.Translator) error {
			return ut.Add("chinacell", "{0}不是一个有效的手机号", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("chinacell", fe.Value().(string))
			return t
		})
		v.validate.RegisterTranslation("email", ZhTrans, func(ut ut.Translator) error {
			return ut.Add("email", "{0}不是一个有效的邮箱", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("email", fe.Value().(string))
			return t
		})
		v.validate.RegisterTranslation("email|chinacell", ZhTrans, func(ut ut.Translator) error {
			return ut.Add("contact", "{0}不是一个有效的联系方式", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("contact", fe.Value().(string))
			return t
		})
		v.validate.RegisterTranslation("password", ZhTrans, func(ut ut.Translator) error {
			return ut.Add("password", "请输入6-20位的密码", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("password")
			return t
		})

		v.validate.RegisterTranslation("required", EnTrans, func(ut ut.Translator) error {
			return ut.Add("required", "{0} is required", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required", fe.Field())
			return t
		})
		v.validate.RegisterTranslation("chinacell", EnTrans, func(ut ut.Translator) error {
			return ut.Add("chinacell", "{0} is not a valid phone number", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("chinacell", fe.Value().(string))
			return t
		})
		v.validate.RegisterTranslation("email", EnTrans, func(ut ut.Translator) error {
			return ut.Add("email", "{0} is not a valid email address", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("email", fe.Value().(string))
			return t
		})
		v.validate.RegisterTranslation("email|chinacell", EnTrans, func(ut ut.Translator) error {
			return ut.Add("contact", "{0} is not a valid contact", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("contact", fe.Value().(string))
			return t
		})
		v.validate.RegisterTranslation("password", EnTrans, func(ut ut.Translator) error {
			return ut.Add("password", "please input a password of 6-20 characters", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("password")
			return t
		})
	})
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
