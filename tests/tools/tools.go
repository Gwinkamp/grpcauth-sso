package tools

import "github.com/brianvoe/gofakeit/v6"

const (
	passDefaultLen = 10
)

func RandomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
