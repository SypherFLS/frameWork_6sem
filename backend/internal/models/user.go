package models

import (
	"framew/internal/log"
)

type User struct {
	Username string
	Password string
}

func (u *User) Validate() *log.Err {
	if u.Username == "" {
		result := log.MakeError(422, "bad user validation", 1, "empty username")
		return result
	} else if u.Password == "" {
		result := log.MakeError(422, "bad user validation", 1, "empty password")
		return result
	}

	result := log.BackSucsess()
	return result
}
