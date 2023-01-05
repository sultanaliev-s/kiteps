package validation_test

import (
	"testing"

	"github.com/sultanaliev-s/kiteps/pkg/validation"
)

type testStruct struct {
	X int    `validate:"required"`
	Y string `validate:"required,min=3,max=10,email"`
}

func TestValidateStruct(t *testing.T) {
	v := validation.NewValidator()

	err := v.ValidateStruct(testStruct{
		X: 0,
		Y: "kiteps@kiteps.kg",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

type testStruct2 struct {
	X int `validate:"test"`
}

func TestRegisterValidation(t *testing.T) {
	v := validation.NewValidator()

	err := v.RegisterValidation("test", func(fl validation.Field) bool {
		return fl.Field().Int() == 1
	}, "must be 1")

	if err != nil {
		t.Fatal("expected nil, got error")
	}

	err = v.ValidateStruct(testStruct2{
		X: 1,
	})
	if err != nil {
		t.Fatal("expected nil, got error")
	}

	err = v.ValidateStruct(testStruct2{
		X: 2,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
