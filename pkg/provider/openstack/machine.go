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

package openstack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alexandrevilain/ollama-machine/pkg/provider"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/images"
)

type MachineManager struct {
	computeClient *gophercloud.ServiceClient
	imageClient   *gophercloud.ServiceClient
}

func newMachineManager(computeClient, imageClient *gophercloud.ServiceClient) *MachineManager {
	return &MachineManager{
		computeClient: computeClient,
		imageClient:   imageClient,
	}
}

func (p *MachineManager) MachineKind() provider.MachineKind {
	return provider.MachineKindVM
}

func (p *MachineManager) Create(ctx context.Context, machineRequest *provider.CreateMachineRequest) (*provider.Machine, error) {
	if err := p.validateCreateMachineRequest(machineRequest); err != nil {
		return nil, err
	}

	flavorID, err := p.findFlavorByName(ctx, machineRequest.InstanceType)
	if err != nil {
		return nil, err
	}

	imageID, err := p.findImageByName(ctx, machineRequest.Image)
	if err != nil {
		return nil, err
	}

	server, err := servers.Create(ctx, p.computeClient, servers.CreateOpts{
		Name:      machineRequest.Name,
		FlavorRef: flavorID,
		ImageRef:  imageID,
		UserData:  machineRequest.UserData,
	}, nil).Extract()
	if err != nil {
		return nil, err
	}

	return serverToMachine(server)
}

func (p *MachineManager) Delete(ctx context.Context, id string) error {
	err := servers.Delete(ctx, p.computeClient, id).ExtractErr()
	if err != nil {
		var unexpectedErr gophercloud.ErrUnexpectedResponseCode
		if errors.As(err, &unexpectedErr) && unexpectedErr.Actual == 404 {
			return nil
		}

		return err
	}

	return nil
}

func (p *MachineManager) Get(ctx context.Context, id string) (*provider.Machine, error) {
	server, err := servers.Get(ctx, p.computeClient, id).Extract()
	if err != nil {
		return nil, err
	}

	return serverToMachine(server)
}

func (p *MachineManager) findFlavorByName(ctx context.Context, name string) (string, error) {
	allPages, err := flavors.ListDetail(p.computeClient, flavors.ListOpts{}).AllPages(ctx)
	if err != nil {
		return "", err
	}

	allFlavors, err := flavors.ExtractFlavors(allPages)
	if err != nil {
		return "", err
	}

	for _, flavor := range allFlavors {
		if flavor.Name == name {
			return flavor.ID, nil
		}
	}

	return "", fmt.Errorf("flavor %q not found", name)
}

func (p *MachineManager) findImageByName(ctx context.Context, name string) (string, error) {
	allPages, err := images.List(p.computeClient, images.ListOpts{}).AllPages(ctx)
	if err != nil {
		return "", err
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		return "", err
	}

	for _, image := range allImages {
		if image.Name == name {
			return image.ID, nil
		}
	}

	return "", fmt.Errorf("image %q not found", name)
}

func (p *MachineManager) validateCreateMachineRequest(req *provider.CreateMachineRequest) error {
	if req.Image == "" {
		return errors.New("image is required with Openstack provider, depending of the cloud provider names can change")
	}

	if req.InstanceType == "" {
		return errors.New("instance type is required")
	}

	return nil
}

func serverToMachine(server *servers.Server) (*provider.Machine, error) {
	state := provider.MachineStatePending
	switch server.Status {
	case "ACTIVE":
		state = provider.MachineStateRunning
	case "ERROR":
		state = provider.MachineStateError
	case "BUILD":
		state = provider.MachineStatePending
	}

	publicIP, err := getPublicIP(server)
	if err != nil {
		return nil, err
	}

	return &provider.Machine{
		ID:    server.ID,
		Name:  server.Name,
		State: state,
		IP:    publicIP,
	}, nil
}

type addresses map[string][]networkInterface

type networkInterface struct {
	Address string `json:"addr"`
	Version int    `json:"version"`
	Type    string `json:"OS-EXT-IPS:type"` //nolint:tagliatelle
}

func getPublicIP(server *servers.Server) (string, error) { //nolint:cyclop
	if server.AccessIPv4 != "" {
		return server.AccessIPv4, nil
	}

	// As gophercloud doesn't provide a struct for server addresses, we need to manually parse it.
	// Got inspiration from CAPO: https://github.com/kubernetes-sigs/cluster-api-provider-openstack/blob/v0.11.4/pkg/cloud/services/compute/instance_types.go#L128
	addressesRaw, err := json.Marshal(server.Addresses)
	if err != nil {
		return "", fmt.Errorf("error marshalling addresses: %w", err)
	}

	var addrs addresses
	err = json.Unmarshal(addressesRaw, &addrs)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling addresses: %w", err)
	}

	// First try to find a floating IP.
	for _, interfaceList := range addrs {
		for _, iface := range interfaceList {
			if iface.Version == 4 && iface.Type == "floating" {
				return iface.Address, nil
			}
		}
	}

	// If no floating IP is found, return the first fixed IP found.
	for _, interfaceList := range addrs {
		for _, iface := range interfaceList {
			if iface.Version == 4 && iface.Type == "fixed" {
				return iface.Address, nil
			}
		}
	}

	return "", nil // Don't return error in this case, as the machine might not have a public IP.
}
