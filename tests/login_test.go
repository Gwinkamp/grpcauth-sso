package tests

import (
	"testing"

	ssov1 "github.com/Gwinkamp/grpcauth-contracts/gen/go"
	"github.com/Gwinkamp/grpcauth-sso/tests/suite"
	"github.com/Gwinkamp/grpcauth-sso/tests/tools"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestLogin_InvalidEmail(t *testing.T) {
	ctx, s := suite.New(t)

	req := &ssov1.LoginRequest{
		Email:     "invalid-email",
		Password:  tools.RandomFakePassword(),
		ServiceId: serviceID,
	}

	resp, err := s.AuthClient.Login(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	assert.Contains(t, err.Error(), "Error:Field validation for 'Email'")
}

func TestLogin_EmptyEmail(t *testing.T) {
	ctx, s := suite.New(t)

	req := &ssov1.LoginRequest{
		Email:     "",
		Password:  tools.RandomFakePassword(),
		ServiceId: serviceID,
	}

	resp, err := s.AuthClient.Login(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	assert.Contains(t, err.Error(), "Error:Field validation for 'Email'")
}

func TestLogin_EmptyPassword(t *testing.T) {
	ctx, s := suite.New(t)

	req := &ssov1.LoginRequest{
		Email:     gofakeit.Email(),
		Password:  "",
		ServiceId: serviceID,
	}

	resp, err := s.AuthClient.Login(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	assert.Contains(t, err.Error(), "Error:Field validation for 'Password'")
}
