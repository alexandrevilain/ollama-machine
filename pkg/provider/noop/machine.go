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

package noop

import (
	"context"

	"github.com/alexandrevilain/ollama-machine/pkg/provider"
)

type MachineManager struct {
	getCount int
}

func newMachineManager() (*MachineManager, error) {
	return &MachineManager{}, nil
}

func (m *MachineManager) MachineKind() provider.MachineKind {
	return provider.MachineKindVM
}

func (m *MachineManager) Create(ctx context.Context, req *provider.CreateMachineRequest) (*provider.Machine, error) {
	return &provider.Machine{
		ID:    "4b00c526-5d3f-4648-b69b-272ab71c6e18",
		Name:  "fake",
		IP:    "1.2.3.4",
		State: provider.MachineStatePending,
	}, nil
}

func (m *MachineManager) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *MachineManager) Get(ctx context.Context, id string) (*provider.Machine, error) {
	if m.getCount < 3 { //nolint:mnd
		m.getCount++

		return &provider.Machine{
			ID:    "4b00c526-5d3f-4648-b69b-272ab71c6e18",
			Name:  "fake",
			IP:    "1.2.3.4",
			State: provider.MachineStatePending,
		}, nil
	}

	return &provider.Machine{
		ID:    "4b00c526-5d3f-4648-b69b-272ab71c6e18",
		Name:  "fake",
		IP:    "1.2.3.4",
		State: provider.MachineStateRunning,
	}, nil
}
