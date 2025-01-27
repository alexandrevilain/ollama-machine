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

package machine_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/alexandrevilain/ollama-machine/pkg/config"
	"github.com/alexandrevilain/ollama-machine/pkg/machine"
	"github.com/alexandrevilain/ollama-machine/pkg/provider"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

func TestSave(t *testing.T) {
	machine.ExportedFS = afero.NewMemMapFs()

	tests := map[string]struct {
		machine *machine.Machine
	}{
		"save machine": {
			machine: &machine.Machine{
				Machine: &provider.Machine{
					ID:    "test-id",
					Name:  "test-name",
					IP:    "127.0.0.1",
					State: provider.MachineStateRunning,
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			err := machine.Save(tt.machine)
			g.Expect(err).NotTo(HaveOccurred())

			filePath := filepath.Join(config.GetMachineDir(), tt.machine.ID+".json")
			file, err := os.Open(filePath)
			g.Expect(err).NotTo(HaveOccurred())
			defer func() {
				_ = file.Close()
			}()

			var savedMachine machine.Machine
			err = json.NewDecoder(file).Decode(&savedMachine)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(savedMachine).To(Equal(*tt.machine))
		})
	}
}

func TestGet(t *testing.T) {
	machine.ExportedFS = afero.NewMemMapFs()

	tests := map[string]struct {
		machine *machine.Machine
	}{
		"get machine": {
			machine: &machine.Machine{
				Machine: &provider.Machine{
					ID:    "test-id",
					Name:  "test-name",
					IP:    "127.0.0.1",
					State: provider.MachineStateRunning,
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			err := machine.Save(tt.machine)
			g.Expect(err).NotTo(HaveOccurred())

			gotMachine, err := machine.Get(tt.machine.ID)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(gotMachine).To(Equal(tt.machine))
		})
	}
}

func TestGetByName(t *testing.T) {
	machine.ExportedFS = afero.NewMemMapFs()

	tests := map[string]struct {
		machine *machine.Machine
	}{
		"get machine by name": {
			machine: &machine.Machine{
				Machine: &provider.Machine{
					ID:    "test-id",
					Name:  "test-name",
					IP:    "127.0.0.1",
					State: provider.MachineStateRunning,
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			err := machine.Save(tt.machine)
			g.Expect(err).NotTo(HaveOccurred())

			gotMachine, err := machine.GetByName(tt.machine.Name)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(gotMachine).To(Equal(tt.machine))
		})
	}
}

func TestList(t *testing.T) {
	machine.ExportedFS = afero.NewMemMapFs()

	tests := map[string]struct {
		machines []*machine.Machine
	}{
		"list machines": {
			machines: []*machine.Machine{
				{
					Machine: &provider.Machine{
						ID:    "test-id-1",
						Name:  "test-name-1",
						IP:    "127.0.0.1",
						State: provider.MachineStateRunning,
					},
				},
				{
					Machine: &provider.Machine{
						ID:    "test-id-2",
						Name:  "test-name-2",
						IP:    "127.0.0.2",
						State: provider.MachineStateStopped,
					},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			for _, m := range tt.machines {
				err := machine.Save(m)
				g.Expect(err).NotTo(HaveOccurred())
			}

			gotMachines, err := machine.List()
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(gotMachines).To(ConsistOf(tt.machines))
		})
	}
}

func TestDelete(t *testing.T) {
	machine.ExportedFS = afero.NewMemMapFs()

	tests := map[string]struct {
		machine *machine.Machine
	}{
		"delete machine": {
			machine: &machine.Machine{
				Machine: &provider.Machine{
					ID:    "test-id",
					Name:  "test-name",
					IP:    "127.0.0.1",
					State: provider.MachineStateRunning,
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			err := machine.Save(tt.machine)
			g.Expect(err).NotTo(HaveOccurred())

			err = machine.Delete(tt.machine.ID)
			g.Expect(err).NotTo(HaveOccurred())

			_, err = machine.Get(tt.machine.ID)
			g.Expect(err).To(HaveOccurred())
		})
	}
}
