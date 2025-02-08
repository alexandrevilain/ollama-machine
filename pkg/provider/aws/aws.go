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

package aws

import (
	"context"
	"errors"

	"github.com/alexandrevilain/ollama-machine/pkg/provider"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
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

func (p *Provider) MachineManager(region string) (provider.MachineManager, error) {
	if p.credentials == nil {
		return nil, errors.New("credentials not set")
	}

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			p.credentials.AccessKeyID,
			p.credentials.SecretAccessKey,
			"", // Session token is empty for long-term credentials
		)),
	)
	if err != nil {
		return nil, err
	}

	client := ec2.NewFromConfig(cfg)

	return newMachineManager(client), nil
}
