package validation

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/mahsandr/arman-challenge/internal/domain/models"
)

type Validator struct {
	validator *validator.Validate
	trans     ut.Translator
}

// NewValidator initializes a new validator.
func NewValidator() *Validator {
	v := &Validator{
		validator: validator.New(),
	}
	enLang := en.New()
	uni := ut.New(enLang, enLang)
	v.trans, _ = uni.GetTranslator("en")
	_ = en_translations.RegisterDefaultTranslations(v.validator, v.trans)
	return v
}

// ValidateStruct validates a struct based on its tags.
func (v *Validator) ValidateStruct(s *models.UserSegment) string {
	err := v.validator.Struct(s)
	if err != nil {
		return err.(validator.ValidationErrors).Error()
	}
	return ""
}
