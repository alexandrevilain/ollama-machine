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
	"errors"

	"github.com/alexandrevilain/ollama-machine/pkg/provider"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
)

type Provider struct {
	credentials *Credentials
}

func NewProvider() *Provider {
	return &Provider{
		credentials: &Credentials{},
	}
}

func (p *Provider) Credentials() provider.Credentials {
	return p.credentials
}

func (p *Provider) MachineManager() (provider.MachineManager, error) {
	if p.credentials == nil {
		return nil, errors.New("credentials not set")
	}

	ctx := context.Background()

	providerClient, err := openstack.AuthenticatedClient(ctx, gophercloud.AuthOptions{
		IdentityEndpoint: p.credentials.IdentityEndpoint,
		Username:         p.credentials.Username,
		Password:         p.credentials.Password,
		TenantID:         p.credentials.TenantID,
		TenantName:       p.credentials.TenantName,
		DomainName:       p.credentials.DomainName,
	})
	if err != nil {
		return nil, err
	}

	computeClient, err := openstack.NewComputeV2(providerClient, gophercloud.EndpointOpts{
		Region: p.credentials.Region,
	})
	if err != nil {
		return nil, err
	}

	imageClient, err := openstack.NewImageV2(providerClient, gophercloud.EndpointOpts{
		Region: p.credentials.Region,
	})
	if err != nil {
		return nil, err
	}

	return newMachineManager(computeClient, imageClient), nil
}
