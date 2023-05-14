package helpers

import (
	"strings"

	"github.com/asaskevich/govalidator"
)

type IValidator interface {
	Validate(someStruct interface{}) ([]string, error)
}

type Validator struct {
}

func NewValidator() IValidator {
	return &Validator{}
}

func (v *Validator) Validate(someStruct interface{}) ([]string, error) {
	_, err := govalidator.ValidateStruct(someStruct)
	if err != nil {
		return strings.Split(err.Error(), ";"), err
	}
	return []string{}, nil
}
