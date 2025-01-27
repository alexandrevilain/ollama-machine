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

package cmd

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/alexandrevilain/ollama-machine/pkg/machine"
	"github.com/alexandrevilain/ollama-machine/pkg/ollama"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	gossh "golang.org/x/crypto/ssh"
)

var (
	localPort  int = ollama.DefaultPort
	remotePort int = ollama.DefaultPort
)

// tunnelCmd represents the tunnel command.
var tunnelCmd = &cobra.Command{
	Use:   "tunnel [machine name]",
	Short: "Create a tunnel to a machine",
	Long: `The tunnel command sets up an SSH tunnel to a specified Ollama machine, enabling secure access to Ollama running on remote machines without exposing it to the internet.

The command creates a local TCP listener that forwards traffic to a remote port on the target machine through an SSH tunnel.

The tunnel command requires exactly one argument:
  - machine name: The name of the machine to establish the tunnel with.

Prerequisites:
  - The target machine must be configured with private connectivity
  - The machine must be running and accessible via SSH on port 22

The tunnel remains active until the process is interrupted (Ctrl+C).`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := machine.GetByName(args[0])
		if err != nil {
			return err
		}

		if m.Connectivity != "private" && m.Connectivity != "" {
			return errors.New("tunneling is only available for machine with private connectivity")
		}

		sshConn, _, err := m.SSHClient()
		if err != nil {
			return fmt.Errorf("failed to create ssh client: %w", err)
		}

		defer func() {
			if closeErr := sshConn.Close(); err != nil {
				log.Error("failed to close ssh connection", "err", closeErr)
			}
		}()

		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", localPort))
		if err != nil {
			return err
		}

		defer func() {
			if closeErr := listener.Close(); err != nil {
				log.Error("failed to close listener connection", "err", closeErr)
			}
		}()

		log.Info("Tunnel available", "localPort", localPort, "remotePort", remotePort)

		for {
			localConn, err := listener.Accept()
			if err != nil {
				log.Fatalf("listen.Accept failed: %v", err)
			}
			go forward(sshConn, localConn, m)
		}
	},
}

func forward(sshClientConn *gossh.Client, localConn net.Conn, m *machine.Machine) {
	sshConn, err := sshClientConn.Dial("tcp", m.OllamaConfig.Address())
	if err != nil {
		log.Fatal("failed to dial ollama", "err", err)
	}

	go func() {
		_, err = io.Copy(sshConn, localConn)
		if err != nil {
			log.Fatal("remote to local data forward failed", "err", err)
		}
	}()

	go func() {
		_, err = io.Copy(localConn, sshConn)
		if err != nil {
			log.Fatal("local to remote data forward failed", "err", err)
		}
	}()
}
