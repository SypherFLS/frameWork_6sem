package models

import (
	"testing"
)

func TestUser_InitUser(t *testing.T) {
	user := &User{}
	result := user.InitUser()

	if result == nil {
		t.Fatal("InitUser вернул nil")
	}

	if result.Username != "" {
		t.Errorf("Ожидалось пустое имя пользователя, получено: %s", result.Username)
	}

	if result.Password != "" {
		t.Errorf("Ожидался пустой пароль, получено: %s", result.Password)
	}

	if result.IsLogin {
		t.Error("IsLogin должен быть false при инициализации")
	}
}

func TestUser_Validate_ValidUser(t *testing.T) {
	user := &User{
		Username: "testuser",
		Password: "testpass",
	}

	err := user.Validate()

	if err == nil {
		t.Fatal("Validate должен вернуть результат")
	}

	if !err.IsDone {
		t.Error("IsDone должен быть true для валидного пользователя")
	}

	if err.Id != 200 {
		t.Errorf("Ожидался код 200, получено: %d", err.Id)
	}
}

func TestUser_Validate_EmptyUsername(t *testing.T) {
	user := &User{
		Username: "",
		Password: "testpass",
	}

	err := user.Validate()

	if err == nil {
		t.Fatal("Validate должен вернуть ошибку")
	}

	if err.IsDone {
		t.Error("IsDone должен быть false для ошибки")
	}

	if err.Id != 422 {
		t.Errorf("Ожидался код 422, получено: %d", err.Id)
	}

	if err.Comment != "bad user validation" {
		t.Errorf("Ожидалось сообщение 'bad user validation', получено: %s", err.Comment)
	}
}

func TestUser_Validate_EmptyPassword(t *testing.T) {
	user := &User{
		Username: "testuser",
		Password: "",
	}

	err := user.Validate()

	if err == nil {
		t.Fatal("Validate должен вернуть ошибку")
	}

	if err.IsDone {
		t.Error("IsDone должен быть false для ошибки")
	}

	if err.Id != 422 {
		t.Errorf("Ожидался код 422, получено: %d", err.Id)
	}

	if err.Comment != "bad user validation" {
		t.Errorf("Ожидалось сообщение 'bad user validation', получено: %s", err.Comment)
	}
}

func TestUser_Validate_BothEmpty(t *testing.T) {
	user := &User{
		Username: "",
		Password: "",
	}

	err := user.Validate()

	if err == nil {
		t.Fatal("Validate должен вернуть ошибку")
	}

	if err.IsDone {
		t.Error("IsDone должен быть false для ошибки")
	}

	if err.Id != 422 {
		t.Errorf("Ожидался код 422, получено: %d", err.Id)
	}

	if err.Identification.Process != "empty username" {
		t.Errorf("Ожидался процесс 'empty username', получено: %s", err.Identification.Process)
	}
}

func TestUser_Login(t *testing.T) {
	user := &User{
		Username: "testuser",
		Password: "testpass",
	}

	err := user.Login()

	if err == nil {
		t.Fatal("Login должен вернуть результат")
	}

	if !err.IsDone {
		t.Error("IsDone должен быть true для успешного входа")
	}

	if err.Id != 201 {
		t.Errorf("Ожидался код 201, получено: %d", err.Id)
	}

	if err.Comment != "login" {
		t.Errorf("Ожидалось сообщение 'login', получено: %s", err.Comment)
	}

	if err.Identification.Process != "user loged in" {
		t.Errorf("Ожидался процесс 'user loged in', получено: %s", err.Identification.Process)
	}
}

func TestUser_Login_MultipleTimes(t *testing.T) {
	user := &User{
		Username: "testuser",
		Password: "testpass",
	}

	err1 := user.Login()
	if err1 == nil || !err1.IsDone {
		t.Fatal("Первый вход должен быть успешным")
	}

	err2 := user.Login()
	if err2 == nil || !err2.IsDone {
		t.Fatal("Второй вход должен быть успешным")
	}

	if err1.Id != err2.Id {
		t.Error("Оба входа должны возвращать одинаковый код")
	}
}

func TestUser_Validate_EdgeCases(t *testing.T) {
	testCases := []struct {
		name     string
		username string
		password string
		expected int
		isDone   bool
	}{
		{
			name:     "Один символ в имени и пароле",
			username: "a",
			password: "b",
			expected: 200,
			isDone:   true,
		},
		{
			name:     "Длинное имя и пароль",
			username: "verylongusernamethatexists",
			password: "verylongpasswordthatexists",
			expected: 200,
			isDone:   true,
		},
		{
			name:     "Пробелы в имени",
			username: "user name",
			password: "password",
			expected: 200,
			isDone:   true,
		},
		{
			name:     "Специальные символы",
			username: "user@123",
			password: "pass!@#",
			expected: 200,
			isDone:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user := &User{
				Username: tc.username,
				Password: tc.password,
			}

			err := user.Validate()

			if err == nil {
				t.Fatal("Validate должен вернуть результат")
			}

			if err.Id != tc.expected {
				t.Errorf("Ожидался код %d, получено: %d", tc.expected, err.Id)
			}

			if err.IsDone != tc.isDone {
				t.Errorf("Ожидался IsDone %v, получено: %v", tc.isDone, err.IsDone)
			}
		})
	}
}

func TestUser_Integration(t *testing.T) {
	user := &User{}
	user = user.InitUser()

	if user.Username != "" || user.Password != "" {
		t.Error("Пользователь должен быть инициализирован с пустыми полями")
	}

	user.Username = "testuser"
	user.Password = "testpass"

	validateErr := user.Validate()
	if validateErr == nil || !validateErr.IsDone {
		t.Fatal("Валидация должна быть успешной")
	}

	loginErr := user.Login()
	if loginErr == nil || !loginErr.IsDone {
		t.Fatal("Вход должен быть успешным")
	}
}
