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

package ssh

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

const (
	DefaultPort = 22
)

// NewClient creates a new SSH client and session.
func NewClient(host, port, user string, k *KeyPairFiles) (*ssh.Client, *ssh.Session, error) {
	privateKeyFile, err := os.ReadFile(k.PrivateKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed reading private key file: %w", err)
	}

	privateKey, err := ssh.ParsePrivateKey(privateKeyFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed parsing private key: %w", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(privateKey)},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey() //nolint:gosec // this should be fixed soon.

	client, err := ssh.Dial("tcp", net.JoinHostPort(host, port), sshConfig)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		_ = client.Close()

		return nil, nil, err
	}

	return client, session, nil
}
