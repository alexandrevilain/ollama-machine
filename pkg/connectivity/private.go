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
	"fmt"

	"github.com/alexandrevilain/ollama-machine/pkg/cloudinit"
	"github.com/alexandrevilain/ollama-machine/pkg/machine"
)

var _ Provider = (*PrivateProvider)(nil)

// PrivateProvider is a private connectivity provider exposing nothing to the outside world.
type PrivateProvider struct{}

// Name returns the name of the provider.
func (p *PrivateProvider) Name() string {
	return "private"
}

// InstallViaCloudInit installs the provider via cloud-init configuration.
// This method does nothing for the private provider.
func (p *PrivateProvider) InstallViaCloudInit(cloudInit *cloudinit.Config) {
	cloudInit.AddRunCmd([]string{"sh", "-c", fmt.Sprintf(`echo "OLLAMA_HOST=localhost" > %s`, machine.OllamaEnvFilePath)})
}

// RetrieveOllamaHost retrieves the Ollama host for the given machine.
// This method always returns "localhost" for the private provider.
func (p *PrivateProvider) RetrieveOllamaHost(m *machine.Machine) (string, error) {
	return "localhost", nil
}
