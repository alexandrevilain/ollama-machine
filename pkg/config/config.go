// Licensed to Alexandre VILAIN under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Alexandre VILAIN licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package config

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

func Init() error {
	return errors.Join(
		os.MkdirAll(GetMachineDir(), 0o750),    //nolint:mnd
		os.MkdirAll(GetMachineKeyDir(), 0o750), //nolint:mnd
	)
}

var baseDir = os.Getenv("OLLAMA_MACHINE_STORAGE_PATH")

// GetMachineDir returns the directory where machines are stored.
func GetMachineDir() string {
	return filepath.Join(getBaseDir(), "machines")
}

// GetMachineKeyDir returns the directory where machine keys are stored.
func GetMachineKeyDir() string {
	return filepath.Join(getBaseDir(), "keys")
}

func getHomeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}

	return os.Getenv("HOME")
}

func getBaseDir() string {
	if baseDir == "" {
		baseDir = filepath.Join(getHomeDir(), ".ollama", "machine")
	}

	return baseDir
}
