package tests

import (
	"testing"

	ssov1 "github.com/Gwinkamp/grpcauth-contracts/gen/go"
	"github.com/Gwinkamp/grpcauth-sso/tests/suite"
	"github.com/Gwinkamp/grpcauth-sso/tests/tools"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegister_InvalidEmail(t *testing.T) {
	ctx, s := suite.New(t)

	resp, err := s.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    "invalid_emaiil",
		Password: tools.RandomFakePassword(),
	})
	require.Error(t, err)

	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	assert.Contains(t, err.Error(), "Error:Field validation for 'Email'")
}

func TestRegister_EmptyEmail(t *testing.T) {
	ctx, s := suite.New(t)

	resp, err := s.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    "",
		Password: tools.RandomFakePassword(),
	})
	require.Error(t, err)

	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	assert.Contains(t, err.Error(), "Error:Field validation for 'Email'")
}

func TestRegister_InvalidShortPassword(t *testing.T) {
	ctx, s := suite.New(t)

	resp, err := s.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    gofakeit.Email(),
		Password: "123",
	})
	require.Error(t, err)

	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	assert.Contains(t, err.Error(), "Error:Field validation for 'Password'")
}

func TestRegister_EmptyPassword(t *testing.T) {
	ctx, s := suite.New(t)

	resp, err := s.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    gofakeit.Email(),
		Password: "",
	})
	require.Error(t, err)

	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	assert.Contains(t, err.Error(), "Error:Field validation for 'Password'")
}

func TestRegister_UserAlreadyRegistered(t *testing.T) {
	ctx, s := suite.New(t)

	email := gofakeit.Email()
	password := tools.RandomFakePassword()

	_, err := s.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)

	resp, err := s.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.Error(t, err)

	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "rpc error: code = AlreadyExists")
	assert.Contains(t, err.Error(), "пользователь уже зарегистрирован")
}
