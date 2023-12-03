package auth

import (
	ssov1 "github.com/Gwinkamp/grpcauth-contracts/gen/go"
	"github.com/go-playground/validator/v10"
)

// registerAuthServerValidationRules задает правила валидации структурам запросов AuthServer
func registerAuthServerValidationRules(validate *validator.Validate) {
	validate.RegisterStructValidationMapRules(
		map[string]string{
			"Email":     "required,email,max=128",
			"Password":  "required,min=4,max=64",
			"ServiceId": "required,uuid4",
		},
		ssov1.LoginRequest{},
	)

	validate.RegisterStructValidationMapRules(
		map[string]string{
			"Email":    "required,email,max=128",
			"Password": "required,min=4,max=64",
		},
		ssov1.RegisterRequest{},
	)

	validate.RegisterStructValidationMapRules(
		map[string]string{
			"UserId": "required,uuid4",
		},
		ssov1.IsAdminRequest{},
	)
}
