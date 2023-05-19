package helpers

import (
	"strings"

	"github.com/asaskevich/govalidator"
)

type IValidator interface {
	Validate(someStruct interface{}) (map[string]interface{}, error)
}

type Validator struct{}

func NewValidator() IValidator {
	return &Validator{}
}

func (v *Validator) Validate(someStruct interface{}) (map[string]interface{}, error) {
	var errorMessages map[string]interface{} = map[string]interface{}{}
	_, err := govalidator.ValidateStruct(someStruct)
	if err != nil {
		for _, msg := range strings.Split(err.Error(), ";") {
			msgDetail := strings.Split(msg, ":")
			key := msgDetail[0]
			val := strings.TrimSpace(msgDetail[1])
			errorMessages[key] = val
		}
		return errorMessages, err
	}
	return errorMessages, nil
}
