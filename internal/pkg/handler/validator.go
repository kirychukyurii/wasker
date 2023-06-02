package handler

import (
	"github.com/go-playground/validator/v10"
)

type Validate struct {
	Validate *validator.Validate
}

type Validator struct {
	validate *validator.Validate
}

func (a *Validator) Validate(i interface{}) error {
	return a.validate.Struct(i)
}
