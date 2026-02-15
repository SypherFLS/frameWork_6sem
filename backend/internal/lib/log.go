package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	logFile   *os.File
	logMutex  sync.Mutex
	logInited bool
)


func InitLogger() error {
	logMutex.Lock()
	defer logMutex.Unlock()

	if logInited {
		return nil
	}


	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("не удалось создать папку logs: %v", err)
		}
	}

	dateStr := time.Now().Format("2006-01-02")
	logFileName := filepath.Join(logDir, fmt.Sprintf("app-%s.log", dateStr))

	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл лога: %v", err)
	}

	logFile = file
	logInited = true
	return nil
}

func CloseLogger() error {
	logMutex.Lock()
	defer logMutex.Unlock()

	if logFile != nil {
		return logFile.Close()
	}
	return nil
}

func writeLog(level, message string) {
	logMutex.Lock()
	defer logMutex.Unlock()

	if !logInited || logFile == nil {
		if err := InitLogger(); err != nil {
			fmt.Printf("[%s] %s: %s\n", time.Now().Format("2006-01-02 15:04:05"), level, message)
			return
		}
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message)

	if _, err := logFile.WriteString(logEntry); err != nil {
		fmt.Printf("Ошибка записи в лог: %v\n", err)
	}

	logFile.Sync()
}

func LogError(err *Err) {
	if err == nil {
		return
	}

	message := fmt.Sprintf("ОШИБКА | Code: %d | Message: %s | RequestID: %d | Process: %s | Success: %v",
		err.Id, err.Comment, err.Identification.Num, err.Identification.Process, err.IsDone)
	writeLog("ERROR", message)
}

func LogOperation(operation string, details string) {
	message := fmt.Sprintf("ОПЕРАЦИЯ | %s | %s", operation, details)
	writeLog("INFO", message)
}

func LogSuccess(err *Err) {
	if err == nil {
		return
	}

	message := fmt.Sprintf("УСПЕХ | Code: %d | Message: %s | RequestID: %d | Process: %s | Success: %v",
		err.Id, err.Comment, err.Identification.Num, err.Identification.Process, err.IsDone)
	writeLog("SUCCESS", message)
}

func Call(id int, com string, num int, proc string) *Err {
	err := MakeError(id, com, num, proc)
	return err
}
