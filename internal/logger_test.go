// Copyright (c) 2015 SUSE LLC. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package containersuseconnect

import (
	"bytes"
	"errors"
	"log"
	"os"
	"strings"
	"testing"
)

func TestGetLogPathDefault(t *testing.T) {
	// Ensure no variable is set
	os.Unsetenv(LogEnv)

	path := GetLogPath()

	if path != DefaultLogPath {
		t.Fatalf("Wrong log file path.\n\tExpected: %s\n\tGot: %s", DefaultLogPath, path)
	}
}

func TestGetLogPathCustom(t *testing.T) {
	if err := os.Setenv(LogEnv, "/tmp/file.log"); err != nil {
		t.Fatalf("Failed to set log file: %v", err)
	}

	file := GetLogPath()
	expected := "/tmp/file.log"

	if file != expected {
		t.Fatalf("Wrong log file path.\n\tExpected: %s\n\tGot: %s", expected, file)
	}
}

// Ensure that the log is always written to a file and Stderr.
func TestSetLoggerOutputToFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "suse.log")

	if err != nil {
		log.Fatalf("Failed to create temp log file: %v", err)
	}

	defer os.Remove(tempFile.Name())

	if err = os.Setenv(LogEnv, tempFile.Name()); err != nil {
		t.Fatalf("Failed to set log file: %v", err)
	}

	defer os.Unsetenv(LogEnv)

	logLine := "This in a log entry in a file and Stderr"

	stdData, err := captureStderr(t, func() {
		SetLoggerOutput()
		log.Println(logLine)
	})

	if err != nil {
		t.Fatalf("Failed to capture Stderr: %v", err)
	}

	buff := new(bytes.Buffer)

	if _, err := buff.ReadFrom(tempFile); err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}

	fileData := buff.String()

	if !strings.Contains(fileData, logLine) {
		t.Fatalf("Wrong log file content. Unable to find the line: %s", logLine)
	}

	if !strings.Contains(stdData, logLine) {
		t.Fatalf("Wrong Stderr content. Unable to find the line: %s", logLine)
	}
}

// Ensure that if the log file is not writable it still writtes to Stderr.
func TestSetLoggerOutputToStderr(t *testing.T) {
	path := "/path/that/does/not/exists/suse.log"

	if err := os.Setenv(LogEnv, path); err != nil {
		t.Fatalf("Failed to set log file: %v", err)
	}

	defer os.Unsetenv(LogEnv)

	logLine := "This in not a log entry in a file"

	data, err := captureStderr(t, func() {
		SetLoggerOutput()
		log.Println(logLine)
	})

	if err != nil {
		t.Fatalf("Failed to capture Stderr: %v", err)
	}

	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("Log file was written")
	}

	if !strings.Contains(data, logLine) {
		t.Fatalf("Wrong Stderr content. Unable to find the line: %s", logLine)
	}
}
