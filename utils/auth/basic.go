package auth

import (
	"encoding/base64"
	"fmt"
)

type Basic struct {
	Name     string
	Password string
}

func NewBasic(name string, password string) *Basic {
	if len(name) == 0 || len(password) == 0 {
		return nil
	}
	return &Basic{
		Name:     name,
		Password: password,
	}
}

func (auth *Basic) FormatHeader() string {
	return "Basic " + base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%v:%v", auth.Name, auth.Password)),
	)
}
