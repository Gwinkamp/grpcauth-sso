package tests

import (
	"testing"
	"time"

	ssov1 "github.com/Gwinkamp/grpcauth-contracts/gen/go"
	"github.com/Gwinkamp/grpcauth-sso/tests/suite"
	"github.com/Gwinkamp/grpcauth-sso/tests/tools"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	tokenSecret = "test-secret"
	serviceID   = "512db16d-6d5b-4af4-aedd-a86e5425df30"
)

func TestRegisterLogin_HappyPath(t *testing.T) {
	ctx, s := suite.New(t)

	email := gofakeit.Email()
	password := tools.RandomFakePassword()

	regResp, err := s.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, regResp.GetUserId())

	defer s.Storage.DeleteUser(ctx, uuid.MustParse(regResp.GetUserId()))

	loginResp, err := s.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:     email,
		Password:  password,
		ServiceId: serviceID,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := loginResp.GetAccessToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, regResp.GetUserId(), claims["uid"])
	assert.Equal(t, email, claims["email"])
	assert.Equal(t, serviceID, claims["service_id"])

	const deltaSeconds = 5

	assert.InDelta(t, loginTime.Add(s.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func TestRegisterLogin_IncorrectPassword(t *testing.T) {
	ctx, s := suite.New(t)

	email := gofakeit.Email()
	password := tools.RandomFakePassword()

	regResp, err := s.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, regResp.GetUserId())

	defer s.Storage.DeleteUser(ctx, uuid.MustParse(regResp.GetUserId()))

	loginResp, err := s.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:     email,
		Password:  "<PASSWORD>",
		ServiceId: serviceID,
	})
	require.Error(t, err)

	assert.Nil(t, loginResp)
	assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	assert.Contains(t, err.Error(), "Неверный логин или пароль")
}

func TestRegisterLogin_IncorrectEmail(t *testing.T) {
	ctx, s := suite.New(t)

	email := gofakeit.Email()
	password := tools.RandomFakePassword()

	regResp, err := s.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, regResp.GetUserId())

	defer s.Storage.DeleteUser(ctx, uuid.MustParse(regResp.GetUserId()))

	loginResp, err := s.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:     "test@test.com",
		Password:  password,
		ServiceId: serviceID,
	})
	require.Error(t, err)

	assert.Nil(t, loginResp)
	assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	assert.Contains(t, err.Error(), "Неверный логин или пароль")
}

func TestRegisterLogin_UndefinedServiceID(t *testing.T) {
	ctx, s := suite.New(t)

	email := gofakeit.Email()
	password := tools.RandomFakePassword()

	regResp, err := s.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, regResp.GetUserId())

	defer s.Storage.DeleteUser(ctx, uuid.MustParse(regResp.GetUserId()))

	loginResp, err := s.AuthClient.Login(ctx, &ssov1.LoginRequest{
        Email:     email,
        Password:  password,
        ServiceId: uuid.NewString(),
    })
	require.Error(t, err)

	assert.Nil(t, loginResp)
	assert.Contains(t, err.Error(), "rpc error: code = InvalidArgument")
	assert.Contains(t, err.Error(), "Сервис не найден")
}
