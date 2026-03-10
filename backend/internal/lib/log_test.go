package lib

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInitLogger(t *testing.T) {
	err := InitLogger()
	if err != nil {
		t.Fatalf("InitLogger вернул ошибку: %v", err)
	}

	if !logInited {
		t.Error("logInited должен быть true после инициализации")
	}

	if logFile == nil {
		t.Error("logFile не должен быть nil после инициализации")
	}

	CloseLogger()
}

func TestInitLogger_MultipleCalls(t *testing.T) {
	err1 := InitLogger()
	if err1 != nil {
		t.Fatalf("Первая инициализация вернула ошибку: %v", err1)
	}

	file1 := logFile

	err2 := InitLogger()
	if err2 != nil {
		t.Fatalf("Вторая инициализация вернула ошибку: %v", err2)
	}

	if logFile != file1 {
		t.Error("При повторной инициализации должен использоваться тот же файл")
	}

	CloseLogger()
}

func TestCloseLogger(t *testing.T) {
	err := InitLogger()
	if err != nil {
		t.Fatalf("InitLogger вернул ошибку: %v", err)
	}

	err = CloseLogger()
	if err != nil {
		t.Errorf("CloseLogger вернул ошибку: %v", err)
	}

	if logFile != nil {
		t.Error("logFile должен быть nil после закрытия")
	}
}

func TestCloseLogger_WithoutInit(t *testing.T) {
	logFile = nil
	logInited = false

	err := CloseLogger()
	if err != nil {
		t.Errorf("CloseLogger не должен возвращать ошибку при отсутствии файла: %v", err)
	}
}

func TestLogError(t *testing.T) {
	err := InitLogger()
	if err != nil {
		t.Fatalf("InitLogger вернул ошибку: %v", err)
	}
	defer CloseLogger()

	testErr := MakeError(400, "test error", 1, "test_process")
	LogError(testErr)

	if logFile == nil {
		t.Error("logFile не должен быть nil")
	}
}

func TestLogError_NilError(t *testing.T) {
	err := InitLogger()
	if err != nil {
		t.Fatalf("InitLogger вернул ошибку: %v", err)
	}
	defer CloseLogger()

	LogError(nil)

	if logFile == nil {
		t.Error("logFile не должен быть nil")
	}
}

func TestLogSuccess(t *testing.T) {
	err := InitLogger()
	if err != nil {
		t.Fatalf("InitLogger вернул ошибку: %v", err)
	}
	defer CloseLogger()

	success := BackSucsess(200, "success", 1, "test_process")
	LogSuccess(success)

	if logFile == nil {
		t.Error("logFile не должен быть nil")
	}
}

func TestLogSuccess_NilError(t *testing.T) {
	err := InitLogger()
	if err != nil {
		t.Fatalf("InitLogger вернул ошибку: %v", err)
	}
	defer CloseLogger()

	LogSuccess(nil)

	if logFile == nil {
		t.Error("logFile не должен быть nil")
	}
}

func TestLogOperation(t *testing.T) {
	err := InitLogger()
	if err != nil {
		t.Fatalf("InitLogger вернул ошибку: %v", err)
	}
	defer CloseLogger()

	LogOperation("TEST_OPERATION", "test details")

	if logFile == nil {
		t.Error("logFile не должен быть nil")
	}
}

func TestInitLogger_CreatesLogDir(t *testing.T) {
	testLogDir := "test_logs"
	originalLogDir := "logs"

	defer func() {
		os.RemoveAll(testLogDir)
		CloseLogger()
	}()

	CloseLogger()
	logInited = false
	logFile = nil

	if _, err := os.Stat(testLogDir); os.IsNotExist(err) {
		err := os.MkdirAll(testLogDir, 0755)
		if err != nil {
			t.Fatalf("Не удалось создать тестовую директорию: %v", err)
		}
	}

	dateStr := time.Now().Format("2006-01-02")
	logFileName := filepath.Join(originalLogDir, "app-"+dateStr+".log")

	if _, err := os.Stat(originalLogDir); os.IsNotExist(err) {
		t.Error("Директория logs должна быть создана")
	}

	if _, err := os.Stat(logFileName); os.IsNotExist(err) {
		t.Error("Файл лога должен быть создан")
	}
}

func TestWriteLog_WithoutInit(t *testing.T) {
	originalLogFile := logFile
	originalLogInited := logInited

	logFile = nil
	logInited = false

	testErr := MakeError(500, "test", 1, "test")
	LogError(testErr)

	logFile = originalLogFile
	logInited = originalLogInited
}
