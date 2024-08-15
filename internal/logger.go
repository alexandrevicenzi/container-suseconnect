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
	"fmt"
	"io"
	"log"
	"os"
)

// The default path for the logger if nothing has been specified.
const DefaultLogPath = "/var/log/suseconnect.log"

// The environment variable used to specify a custom path for the log file.
const LogEnv = "SUSECONNECT_LOG_FILE"

// GetLogPath returns the log file path. If the `LogEnv` environment
// variable has been set, it will return its value if not empty.
// Otherwise, it return `DefaultLogPath`.
func GetLogPath() string {
	if env := os.Getenv(LogEnv); env != "" {
		return env
	}

	return DefaultLogPath
}

func SetLoggerOutput() {
	path := GetLogPath()
	logFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)

	if err == nil {
		writter := io.MultiWriter(os.Stderr, logFile)
		log.SetOutput(writter)
	} else {
		log.SetOutput(os.Stderr)
		log.Printf("Failed to set up log file '%s': %v\n", path, err)
	}
}

// Log the given formatted string with its parameters, and return it
// as a new error.
func loggedError(errorCode int, format string, params ...interface{}) *SuseConnectError {
	msg := fmt.Sprintf(format, params...)
	log.Print(msg)
	return &SuseConnectError{
		ErrorCode: errorCode,
		message:   msg,
	}
}
