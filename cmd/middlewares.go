package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
)

type Middlewares struct{}

type LogEntryModel struct {
	RequestTime  string `json:"request_time"`
	Method       string `json:"method"`
	ResponseTime string `json:"response_time"`
	Status       int    `json:"status"`
	Level        string `json:"level"`
	Error        string `json:"error,omitempty"`
}

func newMiddleware() *Middlewares {
	return &Middlewares{}
}

func (m *Middlewares) LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Start time
		startTime := time.Now()
		err := next(c)

		// Total time taken by the request from start time
		endTime := time.Since(startTime)

		// Creating a log entry with the log entry model
		entry := LogEntryModel{
			RequestTime:  startTime.Format("03:04:05 PM"),
			Method:       c.Request().Method,
			ResponseTime: endTime.String(),
			Status:       c.Response().Status,
		}

		// Checking if the response has any errors if yes then set log level to error else info
		if err != nil {
			entry.Level = "ERROR"
			entry.Error = err.Error()
		} else {
			entry.Level = "INFO"
		}
		writeLog(entry)
		return nil
	}
}

func writeLog(entry LogEntryModel) {
	// Creates a .logs folder which stores every logs in json format
	if err := os.MkdirAll(".logs", 0755); err != nil {
		log.Fatalf("failed to create .logs directory: %v", err)
	}

	// Creating a dynamic file based on date to store all the logs based on the date in sepearte file
	fileName := fmt.Sprintf("%s_logs.log", time.Now().Format("2006-01-02"))
	path := filepath.Join(".logs", fileName)

	// Opening the file to write the log
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open the log file: %v", err)
	}

	// Closing the log at the end of function
	defer file.Close()

	// Encoding to json Format
	data, err := json.Marshal(entry)
	if err != nil {
		log.Fatalf("failed to marshal the log data: %v", err)
	}

	// Writing to json file
	if _, err = file.Write(append(data, '\n')); err != nil {
		log.Fatalf("failed to write the log to json file: %v", err)
	}
}