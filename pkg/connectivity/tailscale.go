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
	"strings"

	"github.com/alexandrevilain/ollama-machine/pkg/cloudinit"
	"github.com/alexandrevilain/ollama-machine/pkg/machine"
)

// TailscaleProvider is a private connectivity provider that exposes ollama through tailscale network.
type TailscaleProvider struct {
	AuthKey string
}

// Name returns the name of the provider.
func (p *TailscaleProvider) Name() string {
	return "tailscale"
}

// InstallViaCloudInit installs tailscale via cloud-init configuration.
func (p *TailscaleProvider) InstallViaCloudInit(cloudInit *cloudinit.Config) {
	cloudInit.AddRunCmd([]string{"sh", "-c", "curl -fsSL https://tailscale.com/install.sh | sh"})
	cloudInit.AddRunCmd([]string{"sh", "-c", "echo 'net.ipv4.ip_forward = 1' | sudo tee -a /etc/sysctl.d/99-tailscale.conf && echo 'net.ipv6.conf.all.forwarding = 1' | sudo tee -a /etc/sysctl.d/99-tailscale.conf && sudo sysctl -p /etc/sysctl.d/99-tailscale.conf"})
	cloudInit.AddRunCmd([]string{"sh", "-c", fmt.Sprintf(`tailscale up --auth-key=%s`, p.AuthKey)})
	cloudInit.AddRunCmd([]string{"sh", "-c", fmt.Sprintf(`echo "OLLAMA_HOST=$(tailscale ip -4)" > %s`, machine.OllamaEnvFilePath)})
}

// RetrieveOllamaHost retrieves the Ollama host for the given machine.
// Under-the-hood is connects to the machine using SSH and runs `tailscale ip -4` to get the machine's IP.
func (p *TailscaleProvider) RetrieveOllamaHost(m *machine.Machine) (string, error) {
	_, sshSession, err := m.SSHClient()
	if err != nil {
		return "", fmt.Errorf("failed to create ssh client: %w", err)
	}

	result, err := sshSession.CombinedOutput("tailscale ip -4")
	if err != nil {
		return "", fmt.Errorf("failed to get machine IP: %w", err)
	}

	ip := strings.Trim(string(result), " \n\t")

	return ip, nil
}
