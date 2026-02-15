package models

import (
	"framew/internal/lib"
	"framew/internal/db"
)

type User struct {
	Username string
	Password string
	IsLogin bool
	storage db.Conteiner
}

func (u *User) InitUser() *User{
	return  &User{

	} 
}

func (u *User) Login() *lib.Err {


	result := lib.BackSucsess(201, "login", 3, "user loged in")
	return result
}


func (u *User) Validate() *lib.Err {
	if u.Username == "" {
		result := lib.MakeError(422, "bad user validation", 1, "empty username")
		return result
	} else if u.Password == "" {
		result := lib.MakeError(422, "bad user validation", 1, "empty password")
		return result
	}

	result := lib.BackSucsess(200, "validation", 1, "user validated")
	return result
}
