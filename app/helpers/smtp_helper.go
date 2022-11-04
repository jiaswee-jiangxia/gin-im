package helpers

import (
	"errors"
	"net/smtp"
)

type LoginAuthWrapper struct { // Our email is unsecured, need to wrap auth
	Username, Password string
}

// LoginAuth is used for smtp login auth
func LoginAuth(username, password string) smtp.Auth {
	return &LoginAuthWrapper{username, password}
}

func (a LoginAuthWrapper) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.Username), nil
}

func (a LoginAuthWrapper) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.Username), nil
		case "Password:":
			return []byte(a.Password), nil
		default:
			return nil, errors.New("Unknown from server")
		}
	}
	return nil, nil
}
