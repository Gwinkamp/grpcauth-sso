package tests

import (
	"testing"

	ssov1 "github.com/Gwinkamp/grpcauth-contracts/gen/go"
	"github.com/Gwinkamp/grpcauth-sso/tests/suite"
	"github.com/Gwinkamp/grpcauth-sso/tests/tools"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegister_InvalidParams(t *testing.T) {
	ctx, s := suite.New(t)

	steps := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "invalid email",
			email:       "invalid-email",
			password:    tools.RandomFakePassword(),
			expectedErr: "Error:Field validation for 'Email'",
		},
		{
			name:        "empty email",
			email:       "",
			password:    tools.RandomFakePassword(),
			expectedErr: "Error:Field validation for 'Email'",
		},
		{
			name:        "empty password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "Error:Field validation for 'Password'",
		},
		{
			name:        "invalid short password",
			email:       gofakeit.Email(),
			password:    "123",
			expectedErr: "Error:Field validation for 'Password'",
		},
	}

	for _, tt := range steps {
		s.Run(tt.name, func(t *testing.T) {
			resp, err := s.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(t, err)

			assert.Nil(t, resp)
			assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestRegister_UserAlreadyRegistered(t *testing.T) {
	ctx, s := suite.New(t)

	email := gofakeit.Email()
	password := tools.RandomFakePassword()

	regResp, err := s.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)

	defer s.Storage.DeleteUser(ctx, uuid.MustParse(regResp.GetUserId()))

	loginResp, err := s.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.Error(t, err)

	assert.Nil(t, loginResp)
	assert.Contains(t, err.Error(), "rpc error: code = AlreadyExists")
	assert.Contains(t, err.Error(), "пользователь уже зарегистрирован")
}
