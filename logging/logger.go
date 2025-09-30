package logging

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
)

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     LogLevel              `json:"level"`
	Message   string                `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

type Logger struct {
	level  LogLevel
	fields map[string]interface{}
}

func NewLogger(level LogLevel) *Logger {
	return &Logger{
		level:  level,
		fields: make(map[string]interface{}),
	}
}

func (l *Logger) WithField(key string, value interface{}) *Logger {
	newLogger := &Logger{
		level:  l.level,
		fields: make(map[string]interface{}),
	}
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}
	newLogger.fields[key] = value
	return newLogger
}

func (l *Logger) log(level LogLevel, message string, fields map[string]interface{}) {
	if l.shouldLog(level) {
		entry := LogEntry{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Level:     level,
			Message:   message,
			Fields:    l.mergeFields(fields),
		}
		
		jsonData, err := json.Marshal(entry)
		if err != nil {
			log.Printf("Failed to marshal log entry: %v", err)
			return
		}
		
		fmt.Fprintln(os.Stdout, string(jsonData))
	}
}

func (l *Logger) shouldLog(level LogLevel) bool {
	levels := map[LogLevel]int{
		DEBUG: 0,
		INFO:  1,
		WARN:  2,
		ERROR: 3,
	}
	return levels[level] >= levels[l.level]
}

func (l *Logger) mergeFields(fields map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})
	for k, v := range l.fields {
		merged[k] = v
	}
	for k, v := range fields {
		merged[k] = v
	}
	return merged
}

func (l *Logger) Debug(message string, fields ...map[string]interface{}) {
	l.log(DEBUG, message, l.mergeFieldsList(fields...))
}

func (l *Logger) Info(message string, fields ...map[string]interface{}) {
	l.log(INFO, message, l.mergeFieldsList(fields...))
}

func (l *Logger) Warn(message string, fields ...map[string]interface{}) {
	l.log(WARN, message, l.mergeFieldsList(fields...))
}

func (l *Logger) Error(message string, fields ...map[string]interface{}) {
	l.log(ERROR, message, l.mergeFieldsList(fields...))
}

func (l *Logger) mergeFieldsList(fieldsList ...map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})
	for _, fields := range fieldsList {
		for k, v := range fields {
			merged[k] = v
		}
	}
	return merged
}
