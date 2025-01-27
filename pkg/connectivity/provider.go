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

package connectivity

import (
	"github.com/alexandrevilain/ollama-machine/pkg/cloudinit"
	"github.com/alexandrevilain/ollama-machine/pkg/machine"
)

// Provider is an interface that defines the methods that a connectivity provider must implement.
type Provider interface {
	Name() string
	InstallViaCloudInit(cloudInit *cloudinit.Config)
	RetrieveOllamaHost(m *machine.Machine) (string, error)
}

// Options are the options for a machine connectivity.
type Options struct {
	Public           bool
	TailscaleAuthKey string
}

// GetProvider returns the appropriate provider based on the options.
func GetProvider(opts *Options) Provider {
	if opts.Public {
		return &PublicProvider{}
	}

	if opts.TailscaleAuthKey != "" {
		return &TailscaleProvider{
			AuthKey: opts.TailscaleAuthKey,
		}
	}

	return &PrivateProvider{}
}
