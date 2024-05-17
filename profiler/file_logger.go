package profiler

import (
	"fmt"
	"os"
	"time"
)

type FileLoggerData struct {
	data string
}

type FileLogger struct {
	filename string
	file     *os.File
	data     chan FileLoggerData
	stop     chan struct{}
}

func NewFileLogger(filename string) *FileLogger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error creating log file: %v\n", err)
		return nil
	}

	return &FileLogger{
		filename: filename,
		file:     logFile,
		data:     make(chan FileLoggerData),
		stop:     make(chan struct{}),
	}
}

func (fl *FileLogger) Write(d FileLoggerData) {
	_, err := fmt.Fprintf(fl.file, "%v;%s\n", time.Now().Unix(), d)
	if err != nil {
		fmt.Printf("Error writing to log file '%s': %v\n", fl.filename, err)
	}
}

func (fl *FileLogger) Listen() {
	go func() {
		for d := range fl.data {
			fl.Write(d)
		}
	}()
}

func (fl *FileLogger) Close() {
	fl.file.Close()
	close(fl.data)
	close(fl.stop)
}
