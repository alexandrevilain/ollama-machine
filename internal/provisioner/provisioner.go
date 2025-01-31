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

package provisioner

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"syscall"
	"time"

	"github.com/alexandrevilain/ollama-machine/pkg/cloudcredentials"
	"github.com/alexandrevilain/ollama-machine/pkg/cloudinit"
	"github.com/alexandrevilain/ollama-machine/pkg/config"
	"github.com/alexandrevilain/ollama-machine/pkg/connectivity"
	"github.com/alexandrevilain/ollama-machine/pkg/machine"
	"github.com/alexandrevilain/ollama-machine/pkg/ollama"
	"github.com/alexandrevilain/ollama-machine/pkg/provider"
	"github.com/alexandrevilain/ollama-machine/pkg/registry"
	"github.com/alexandrevilain/ollama-machine/pkg/ssh"
	"github.com/charmbracelet/log"
)

var waitMachineStateInterval = 5 * time.Second

// Provisioner is responsible for creating and deleting machines.
type Provisioner struct {
	providerName    string
	credentialsName string
	machineManager  provider.MachineManager
}

// NewProvisioner creates a new instance of provisioner.
func NewProvisioner(providerName, credentialsName string) (*Provisioner, error) {
	provider, ok := registry.Providers[providerName]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", providerName)
	}

	err := cloudcredentials.Get(cloudcredentials.Key{
		Name:     credentialsName,
		Provider: providerName,
	}, provider.Credentials())
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	machineManager, err := provider.MachineManager()
	if err != nil {
		return nil, fmt.Errorf("failed to get machine manager: %w", err)
	}

	return &Provisioner{
		providerName:    providerName,
		credentialsName: credentialsName,
		machineManager:  machineManager,
	}, nil
}

// NewProvisionerForMachine creates a new instance of provisioner for a specific machine.
func NewProvisionerForMachine(machineName string) (*Provisioner, error) {
	m, err := machine.GetByName(machineName)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine: %w", err)
	}

	providerName := m.ProviderName
	credentialsName := m.CredentialsName

	provider, ok := registry.Providers[providerName]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", providerName)
	}

	err = cloudcredentials.Get(cloudcredentials.Key{
		Name:     credentialsName,
		Provider: providerName,
	}, provider.Credentials())
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	machineManager, err := provider.MachineManager()
	if err != nil {
		return nil, fmt.Errorf("failed to get machine manager: %w", err)
	}

	return &Provisioner{
		providerName:    providerName,
		credentialsName: credentialsName,
		machineManager:  machineManager,
	}, nil
}

// CreateMachine creates a new machine.
func (p *Provisioner) CreateMachine(ctx context.Context, req *provider.CreateMachineRequest, connectivityOpts *connectivity.Options) error { //nolint:funlen,cyclop
	connectivityProvider := connectivity.GetProvider(connectivityOpts)

	log.Info("Generating SSH key pair")

	keyPair, keyPairFiles, err := ssh.GenerateSSHKey(config.GetMachineKeyDir(), req.Name)
	if err != nil {
		return err
	}

	log.Info("Generating machine config")

	if p.machineManager.MachineKind() == provider.MachineKindVM {
		cloudInit := p.generateCloudInit(connectivityProvider, keyPair) // TODO(alexandrevilain): this is a great v0 but it should be improved.

		req.UserData, err = cloudInit.Render()
		if err != nil {
			return err
		}
	}

	log.Info("Creating machine")
	providerMachine, err := p.machineManager.Create(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create machine: %w", err)
	}

	m := &machine.Machine{
		Machine:         providerMachine,
		ProviderName:    p.providerName,
		CredentialsName: p.credentialsName,
		KeyPair:         keyPairFiles,
		Connectivity:    connectivityProvider.Name(),
	}

	// Start by saving the machine before waiting for it to be ready
	log.Info("Saving machine configuration to disk")
	err = machine.Save(m)
	if err != nil {
		return fmt.Errorf("failed to save machine: %w", err)
	}

	log.Info("Waiting for machine to be ready")
	for {
		providerMachine, err = p.machineManager.Get(ctx, providerMachine.ID)
		if err != nil {
			return fmt.Errorf("failed to get machine: %w", err)
		}

		if providerMachine.State == provider.MachineStateError {
			return fmt.Errorf("machine %s is in error state", providerMachine.ID)
		}

		if providerMachine.State == provider.MachineStateRunning {
			m.Machine = providerMachine

			log.Info("Machine ready")

			break
		}

		time.Sleep(waitMachineStateInterval)

		log.Info("Still waiting for machine to be ready")
	}

	// Machine is ready save it.
	m.Machine = providerMachine
	err = machine.Save(m)
	if err != nil {
		return fmt.Errorf("failed to save machine: %w", err)
	}

	log.Info("Waiting for Ollama to be started")

	// TODO(alexandrevilain): This is a temporary solution to wait for Ollama to be started.
	// We should have a better way to know when Ollama is ready.
	// This would not work if we're asking Ollama to pre-pull models for instance.
	for {
		_, sshSession, err := m.SSHClient()
		if err != nil {
			if errors.Is(err, syscall.ECONNREFUSED) {
				log.Info("Waiting for SSH to be ready", "err", err)

				time.Sleep(waitMachineStateInterval)

				continue
			}

			return fmt.Errorf("failed to create ssh client: %w", err)
		}

		result, err := sshSession.CombinedOutput("systemctl is-active ollama")
		if err != nil {
			log.Info("Still waiting for Ollama to be started", "err", err)
		}

		status := strings.Trim(string(result), " \n\t")

		if status == "active" {
			log.Info("Ollama started")

			break
		}

		log.Info("Still waiting for Ollama to be started", "status", status)

		time.Sleep(waitMachineStateInterval)
	}

	log.Info("Retrieving Ollama host")
	m.OllamaConfig.Host, err = connectivityProvider.RetrieveOllamaHost(m)
	if err != nil {
		return fmt.Errorf("failed to retrieve Ollama host IP from connectivity provider: %w", err)
	}
	m.OllamaConfig.Port = ollama.DefaultPort

	err = machine.Save(m)
	if err != nil {
		return fmt.Errorf("failed to save machine: %w", err)
	}

	log.Info("Machine ready!")

	return nil
}

func (p *Provisioner) generateCloudInit(connectivityProvider connectivity.Provider, keyPair *ssh.KeyPair) *cloudinit.Config {
	cloudInit := cloudinit.NewConfig()
	cloudInit.AddUser(cloudinit.User{
		Name:   machine.SSHUsername,
		Groups: "sudo",
		Shell:  "/bin/bash",
		Sudo:   "ALL=(ALL) NOPASSWD:ALL",
		SSHAuthorizedKeys: []string{
			string(keyPair.PublicKey),
		},
	})

	connectivityProvider.InstallViaCloudInit(cloudInit)

	cloudInit.AddFile(cloudinit.File{
		Path: "/etc/systemd/system/ollama.service.d/override.conf",
		Content: fmt.Sprintf(`[Service]
EnvironmentFile=%s`, machine.OllamaEnvFilePath),
	})
	cloudInit.AddRunCmd([]string{"sh", "-c", "curl -fsSL https://ollama.com/install.sh | sh"})
	cloudInit.AddRunCmd([]string{"sh", "-c", "sudo systemctl start ollama"})

	return cloudInit
}

// DeleteMachine deletes a machine.
func (p *Provisioner) DeleteMachine(ctx context.Context, machineName string) error {
	m, err := machine.GetByName(machineName)
	if err != nil {
		return fmt.Errorf("failed to get machine: %w", err)
	}

	log.Info("Deleting machine")

	err = p.machineManager.Delete(ctx, m.ID)
	if err != nil {
		return fmt.Errorf("failed to delete machine: %w", err)
	}

	if m.KeyPair != nil {
		log.Info("Deleting key pair files")

		if err := ssh.DeleteKeyPairFiles(m.KeyPair); err != nil {
			return fmt.Errorf("failed to delete key pair files: %w", err)
		}
	}

	log.Info("Deleting machine configuration")

	err = machine.Delete(m.ID)
	if err != nil {
		return fmt.Errorf("failed to delete machine: %w", err)
	}

	log.Info("Machine deleted")

	return nil
}

func (p *Provisioner) StartMachine(ctx context.Context, machineName string) error {
	m, err := machine.GetByName(machineName)
	if err != nil {
		return fmt.Errorf("failed to get machine: %w", err)
	}

	log.Info("Starting machine")

	err = p.machineManager.Start(ctx, m.ID)
	if err != nil {
		return fmt.Errorf("failed to start machine: %w", err)
	}

	for {
		providerMachine, err := p.machineManager.Get(ctx, m.ID)
		if err != nil {
			return fmt.Errorf("failed to get machine: %w", err)
		}

		if providerMachine.State == provider.MachineStateRunning {
			m.Machine = providerMachine

			break
		}

		log.Info("Still waiting for machine to be started")
		time.Sleep(waitMachineStateInterval)
	}

	log.Info("Machine started")

	err = machine.Save(m)
	if err != nil {
		return fmt.Errorf("failed to save machine: %w", err)
	}

	return nil
}

func (p *Provisioner) StopMachine(ctx context.Context, machineName string) error {
	m, err := machine.GetByName(machineName)
	if err != nil {
		return fmt.Errorf("failed to get machine: %w", err)
	}

	log.Info("Stopping machine")

	err = p.machineManager.Stop(ctx, m.ID)
	if err != nil {
		return fmt.Errorf("failed to stop machine: %w", err)
	}

	for {
		providerMachine, err := p.machineManager.Get(ctx, m.ID)
		if err != nil {
			return fmt.Errorf("failed to get machine: %w", err)
		}

		if providerMachine.State == provider.MachineStateStopped {
			m.Machine = providerMachine

			break
		}

		log.Info("Still waiting for machine to be stopped")
		time.Sleep(waitMachineStateInterval)
	}

	err = machine.Save(m)
	if err != nil {
		return fmt.Errorf("failed to save machine: %w", err)
	}

	log.Info("Machine stopped")

	return nil
}
