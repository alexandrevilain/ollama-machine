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

package provider

import (
	"context"

	"github.com/spf13/pflag"
)

// Credentials represents the interface for managing provider credentials.
type Credentials interface {
	// RegisterFlags registers the necessary flags for credentials.
	RegisterFlags(fs *pflag.FlagSet)
	// Validate validates the credentials.
	Validate() error
	// Complete completes the credentials.
	// This can be useful for reading a password from stdin.
	Complete() error
}

// MachineManager represents the interface for managing machines.
type MachineManager interface {
	// Create creates a new machine with the given request.
	Create(ctx context.Context, machine *CreateMachineRequest) (*Machine, error)
	// Delete deletes the machine with the given ID.
	Delete(ctx context.Context, id string) error
	// Start starts the machine with the given ID.
	Start(ctx context.Context, id string) error
	// Stop stops the machine with the given ID.
	// This is a soft stop, meaning the machine can be started again.
	// But underlying providers should ensure the stop implementation
	// results in the machine not incurring too costs (snapshot or storage cost only).
	Stop(ctx context.Context, id string) error
	// Get retrieves the machine with the given ID.
	Get(ctx context.Context, id string) (*Machine, error)
	// MachineKind returns the kind of machine managed.
	MachineKind() MachineKind
}

// Provider represents the interface for a cloud provider.
type Provider interface {
	// Credentials returns the credentials for the provider.
	Credentials() Credentials
	// MachineManager returns the machine manager for the provider.
	MachineManager(region string) (MachineManager, error)
}

// Machine represents a machine instance.
type Machine struct {
	// ID is the unique identifier of the machine.
	ID string `json:"id"`
	// Name is the name of the machine.
	Name string `json:"name"`
	// IP is the IP address of the machine.
	IP string `json:"ip"`
	// Region is the region the machine is in.
	Region string `json:"region"`
	// State is the current state of the machine.
	State MachineState `json:"state"`
}

// CreateMachineRequest represents the request to create a new machine.
type CreateMachineRequest struct {
	// Name is the name of the machine.
	Name string
	// InstanceType is the type of instance to create.
	InstanceType string
	// Image is the image to use for the machine.
	Image string
	// Zone is the zone to create the machine in.
	Zone string
	// Tags are the tags to assign to the machine.
	Tags map[string]string
	// UserData is the user data to provide to the machine.
	UserData []byte
}

// MachineState represents the state of a machine.
type MachineState string

const (
	// MachineStatePending indicates the machine is pending.
	MachineStatePending MachineState = "pending"
	// MachineStateRunning indicates the machine is running.
	MachineStateRunning MachineState = "running"
	// MachineStateStopped indicates the machine is stopped.
	MachineStateStopped MachineState = "stopped"
	// MachineStateTerminated indicates the machine is terminated.
	MachineStateTerminated MachineState = "terminated"
	// MachineStateError indicates the machine is in an error state.
	MachineStateError MachineState = "error"
)

// MachineKind represents the kind of machine.
type MachineKind string

const (
	// MachineKindVM indicates the machine is a virtual machine.
	MachineKindVM MachineKind = "vm"
	// MachineKindContainer indicates the machine is a container.
	MachineKindContainer MachineKind = "container"
)
