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

package machine

import (
	"net"
	"strconv"

	"github.com/alexandrevilain/ollama-machine/pkg/provider"
	"github.com/alexandrevilain/ollama-machine/pkg/ssh"
	gossh "golang.org/x/crypto/ssh"
)

const (
	SSHUsername       = "ollama-machine"
	OllamaEnvFilePath = "/home/ollama-machine/env"
)

type Machine struct {
	*provider.Machine

	OllamaConfig    OllamaConfig      `json:"ollamaConfig"`
	ProviderName    string            `json:"providerName"`
	CredentialsName string            `json:"credentialsName"`
	Connectivity    string            `json:"connectivity"`
	KeyPair         *ssh.KeyPairFiles `json:"keyPair"`
}

// SSHClient returns a new SSH client and session for the machine.
func (m *Machine) SSHClient() (*gossh.Client, *gossh.Session, error) {
	return ssh.NewClient(m.IP, "22", SSHUsername, m.KeyPair)
}

type OllamaConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (o OllamaConfig) Address() string {
	return net.JoinHostPort(o.Host, strconv.Itoa(o.Port))
}
