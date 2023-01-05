package validation

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	vLib "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type (
	Validator struct {
		validation *vLib.Validate
		trans      ut.Translator
	}

	// Field describes a field in a struct.
	// It is used to get info about it for validation.
	Field vLib.FieldLevel
)

// Our validator is setup only for 'en' locale.
// If we need to support other locales, we can add them here.
// But do not forget to add locales to the RegisterValidation too then.
func NewValidator() *Validator {
	uni := ut.New(en.New())
	trans, _ := uni.GetTranslator("en")
	v := vLib.New()
	en_translations.RegisterDefaultTranslations(v, trans)

	return &Validator{
		validation: v,
		trans:      trans,
	}
}

// RegisterValidation registers a new validation function with the validator.
// This function will be called when the validator encounters the tag.
func (v *Validator) RegisterValidation(tag string, check func(Field) bool, errMsg string) error {
	err := v.validation.RegisterTranslation(tag, v.trans,
		func(ut ut.Translator) error {
			return ut.Add(tag, errMsg, true)
		},
		func(ut ut.Translator, fe vLib.FieldError) string {
			t, err := ut.T(tag, fe.Field())
			if err != nil {
				// TODO: change this to something more error tolerant
				panic("could not register validation")
			}
			return t
		},
	)
	if err != nil {
		return err
	}
	return v.validation.RegisterValidation(tag, func(fl vLib.FieldLevel) bool {
		return check(fl)
	})
}

// ValidateStruct validates the given struct.
// It validates it according to the tags on the struct.
// You can lookup available tags [here].
// Or create custom tags by using RegisterValidation.
//
// [here]: https://github.com/go-playground/validator
func (v *Validator) ValidateStruct(value any) error {
	return v.validation.Struct(value)
}

// UnpackErrors unpacks the error returned by ValidateStruct into a slice of strings.
func (v *Validator) UnpackErrors(e error) []string {
	values, ok := e.(vLib.ValidationErrors)
	if !ok {
		return nil
	}
	errs := make([]string, 0, len(values))
	for _, vv := range values {
		errs = append(errs, vv.Translate(v.trans))
	}
	return errs
}
