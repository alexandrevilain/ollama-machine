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

package machine

import (
	"encoding/json"
	"fmt"
	iofs "io/fs"
	"path/filepath"

	"github.com/alexandrevilain/ollama-machine/pkg/config"
	"github.com/spf13/afero"
)

// TODO(alexandrevilain): once we have a common way to inject dependencies to commands,
// we should use a Store object with the filesystem instead of using a global variable.
var fs = afero.NewOsFs() //nolint:varnamelen

// Save saves the given machine to a file.
func Save(machine *Machine) error {
	file, err := fs.Create(machineFilename(machine.ID))
	if err != nil {
		return err
	}

	return json.NewEncoder(file).Encode(machine)
}

// Get retrieves a machine by its ID.
func Get(id string) (*Machine, error) {
	file, err := fs.Open(machineFilename(id))
	if err != nil {
		return nil, err
	}

	result := &Machine{}
	err = json.NewDecoder(file).Decode(result)

	return result, err
}

// GetByName retrieves a machine by its name.
func GetByName(name string) (*Machine, error) {
	machines, err := List()
	if err != nil {
		return nil, err
	}

	for _, machine := range machines {
		if machine.Name == name {
			return machine, nil
		}
	}

	return nil, fmt.Errorf("machine %s not found", name)
}

// List lists all machines.
func List() ([]*Machine, error) {
	result := []*Machine{}
	err := afero.Walk(fs, config.GetMachineDir(), func(path string, info iofs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		file, err := fs.Open(path)
		if err != nil {
			return err
		}

		machine := &Machine{}
		err = json.NewDecoder(file).Decode(machine)
		if err != nil {
			return err
		}

		result = append(result, machine)

		return nil
	})

	return result, err
}

// Delete deletes a machine by its ID.
func Delete(id string) error {
	return fs.Remove(machineFilename(id))
}

// machineFilename returns the filename for the machine with the given ID.
func machineFilename(id string) string {
	return filepath.Join(config.GetMachineDir(), id+".json")
}
