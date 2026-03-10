package lib

import (
	"testing"
)

func TestMakeError(t *testing.T) {
	err := MakeError(400, "bad request", 1, "test_process")

	if err == nil {
		t.Fatal("MakeError вернул nil")
	}

	if err.Id != 400 {
		t.Errorf("Ожидался код ошибки 400, получено: %d", err.Id)
	}

	if err.Comment != "bad request" {
		t.Errorf("Ожидалось сообщение 'bad request', получено: %s", err.Comment)
	}

	if err.IsDone != false {
		t.Error("IsDone должен быть false для ошибки")
	}

	if err.Identification.Num != 1 {
		t.Errorf("Ожидался номер идентификации 1, получено: %d", err.Identification.Num)
	}

	if err.Identification.Process != "test_process" {
		t.Errorf("Ожидался процесс 'test_process', получено: %s", err.Identification.Process)
	}
}

func TestBackSucsess(t *testing.T) {
	err := BackSucsess(200, "success", 2, "test_process")

	if err == nil {
		t.Fatal("BackSucsess вернул nil")
	}

	if err.Id != 200 {
		t.Errorf("Ожидался код 200, получено: %d", err.Id)
	}

	if err.Comment != "success" {
		t.Errorf("Ожидалось сообщение 'success', получено: %s", err.Comment)
	}

	if err.IsDone != true {
		t.Error("IsDone должен быть true для успешного результата")
	}

	if err.Identification.Num != 2 {
		t.Errorf("Ожидался номер идентификации 2, получено: %d", err.Identification.Num)
	}

	if err.Identification.Process != "test_process" {
		t.Errorf("Ожидался процесс 'test_process', получено: %s", err.Identification.Process)
	}
}

func TestMakeError_DifferentCodes(t *testing.T) {
	testCases := []struct {
		code    int
		message string
	}{
		{400, "bad request"},
		{404, "not found"},
		{422, "validation error"},
		{500, "internal server error"},
	}

	for _, tc := range testCases {
		err := MakeError(tc.code, tc.message, 1, "test")

		if err.Id != tc.code {
			t.Errorf("Ожидался код %d, получено: %d", tc.code, err.Id)
		}

		if err.Comment != tc.message {
			t.Errorf("Ожидалось сообщение '%s', получено: %s", tc.message, err.Comment)
		}

		if err.IsDone {
			t.Error("IsDone должен быть false для ошибки")
		}
	}
}

func TestBackSucsess_DifferentCodes(t *testing.T) {
	testCases := []struct {
		code    int
		message string
	}{
		{200, "ok"},
		{201, "created"},
		{204, "no content"},
	}

	for _, tc := range testCases {
		err := BackSucsess(tc.code, tc.message, 1, "test")

		if err.Id != tc.code {
			t.Errorf("Ожидался код %d, получено: %d", tc.code, err.Id)
		}

		if err.Comment != tc.message {
			t.Errorf("Ожидалось сообщение '%s', получено: %s", tc.message, err.Comment)
		}

		if !err.IsDone {
			t.Error("IsDone должен быть true для успешного результата")
		}
	}
}

func TestErr_Identification(t *testing.T) {
	err := MakeError(400, "test", 123, "my_process")

	if err.Identification.Num != 123 {
		t.Errorf("Ожидался номер 123, получено: %d", err.Identification.Num)
	}

	if err.Identification.Process != "my_process" {
		t.Errorf("Ожидался процесс 'my_process', получено: %s", err.Identification.Process)
	}
}
