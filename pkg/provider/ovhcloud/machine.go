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

package ovhcloud

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/alexandrevilain/ollama-machine/pkg/provider"
	ovhsdk "github.com/dirien/ovh-go-sdk/pkg/sdk"
	"github.com/ovh/go-ovh/ovh"
)

type MachineManager struct {
	client *ovhsdk.OVHcloud
}

func newMachineManager(client *ovhsdk.OVHcloud) *MachineManager {
	return &MachineManager{client: client}
}

func (m *MachineManager) MachineKind() provider.MachineKind {
	return provider.MachineKindVM
}

func (m *MachineManager) Create(ctx context.Context, req *provider.CreateMachineRequest) (*provider.Machine, error) {
	if err := m.validateCreateMachineRequest(req); err != nil {
		return nil, err
	}

	image, err := m.client.GetImage(ctx, req.Image, req.Region)
	if err != nil {
		return nil, err
	}

	flavor, err := m.client.GetFlavor(ctx, req.InstanceType, req.Region)
	if err != nil {
		return nil, err
	}

	instance, err := m.client.CreateInstance(ctx, ovhsdk.InstanceCreateOptions{
		Name:           req.Name,
		Region:         req.Region,
		FlavorID:       flavor.ID,
		ImageID:        image.ID,
		MonthlyBilling: false,
		UserData:       string(req.UserData),
	})
	if err != nil {
		return nil, err
	}

	return instanceToMachine(instance)
}

func (m *MachineManager) Delete(ctx context.Context, id string) error {
	err := m.client.DeleteInstance(ctx, id)
	if err != nil {
		var apiError *ovh.APIError
		if errors.As(err, &apiError) {
			if apiError.Code == http.StatusNotFound {
				return nil
			}
		}

		return err
	}

	return nil
}

func (m *MachineManager) Start(ctx context.Context, id string) error {
	url := fmt.Sprintf("/cloud/project/%s/instance/%s/unshelve", m.client.ServiceName, id)

	return m.client.Client.PostWithContext(ctx, url, nil, nil)
}

func (m *MachineManager) Stop(ctx context.Context, id string) error {
	url := fmt.Sprintf("/cloud/project/%s/instance/%s/shelve", m.client.ServiceName, id)

	return m.client.Client.PostWithContext(ctx, url, nil, nil)
}

func (m *MachineManager) Get(ctx context.Context, id string) (*provider.Machine, error) {
	instance, err := m.client.GetInstance(ctx, id)
	if err != nil {
		return nil, err
	}

	return instanceToMachine(instance)
}

func (m *MachineManager) validateCreateMachineRequest(req *provider.CreateMachineRequest) error {
	if req.InstanceType == "" {
		return errors.New("instance type is required")
	}

	if req.Image == "" {
		req.Image = "ubuntu-24.04"
	}

	if req.Region == "" {
		req.Region = "GRA7"
	}

	return nil
}

func instanceToMachine(instance *ovhsdk.Instance) (*provider.Machine, error) {
	ipv4, _ := ovhsdk.IPv4(instance)

	var state provider.MachineState
	switch instance.Status {
	case ovhsdk.InstanceActive:
		state = provider.MachineStateRunning
	case ovhsdk.InstanceDeleted:
		state = provider.MachineStateTerminated
	case ovhsdk.InstanceError:
		state = provider.MachineStateError
	case ovhsdk.InstanceStopped,
		ovhsdk.InstanceStatus("SHELVED"),
		ovhsdk.InstanceStatus("SHELVED_OFFLOADED"):
		state = provider.MachineStateStopped
	case ovhsdk.InstanceDeleting,
		ovhsdk.InstanceReboot,
		ovhsdk.InstanceBuilding,
		ovhsdk.InstanceUnknown,
		ovhsdk.InstanceBuild,
		ovhsdk.InstanceResuming,
		ovhsdk.InstanceRebuild,
		ovhsdk.InstanceStatus("UNSHELVING"):
		state = provider.MachineStatePending
	}

	return &provider.Machine{
		ID:    instance.ID,
		Name:  instance.Name,
		IP:    ipv4,
		State: state,
	}, nil
}
