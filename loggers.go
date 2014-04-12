// Package loggers implements routines to log errors and access
// for a server.
package loggers

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Error creates an error logger to the given file.
func Error(file string) (*log.Logger, *os.File, error) {
	err := os.MkdirAll(filepath.Dir(file), os.ModePerm|os.ModeDir)
	if err != nil {
		return nil, nil, err
	}

	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, nil, err
	}

	return log.New(logFile, "", log.LstdFlags), logFile, nil
}

// Access creates or opens a log file in dir. A new log file is only created
// if there's no log file less than a week old.
func Access(dir string) (*os.File, error) {
	var logFile *os.File

	err := os.MkdirAll(dir, os.ModePerm|os.ModeDir)
	if err != nil {
		return nil, err
	}

	logDir, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer logDir.Close()

	logFiles, err := logDir.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	// Find the most recent if files exist.
	if len(logFiles) > 0 {
		for _, name := range logFiles {
			base := strings.Replace(name, filepath.Ext(name), "", 1)

			fileDate, err := time.Parse(time.RFC3339, base)
			if err != nil {
				continue
			}

			// Use the file if less than a week old.
			if time.Since(fileDate).Hours() <= 168 {
				path := filepath.Join(dir, name)
				logFile, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0644)
				if err != nil {
					return nil, err
				}
				break
			}
		}
	}

	// Create new file if none have been found.
	if logFile == nil {
		path := filepath.Join(dir, time.Now().Format(time.RFC3339)+".log")
		logFile, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}
	}

	return logFile, nil
}
