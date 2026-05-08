package validator

import (
	validatorPkg "github.com/go-playground/validator/v10"
)

type Validator interface {
	ValidateJSON(any) error
}

func New() Validator {
	return validatorService{
		validate: validatorPkg.New(validatorPkg.WithRequiredStructEnabled()),
	}
}

type validatorService struct {
	validate *validatorPkg.Validate
}

func (v validatorService) ValidateJSON(data any) error {
	return v.validate.Struct(data)
}
