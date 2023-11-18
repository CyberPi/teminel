package auth

import (
	"encoding/base64"
	"fmt"
)

type Basic struct {
	Name     string
	Password string
}

func (auth *Basic) FormatHeader() string {
	return "Basic " + base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%v:%v", auth.Name, auth.Password)),
	)
}
