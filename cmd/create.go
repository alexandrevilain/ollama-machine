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
	"fmt"

	"github.com/alexandrevilain/ollama-machine/internal/provisioner"
	"github.com/alexandrevilain/ollama-machine/pkg/connectivity"
	"github.com/alexandrevilain/ollama-machine/pkg/provider"
	"github.com/spf13/cobra"
)

var (
	createRequest    = &provider.CreateMachineRequest{}
	connectivityOpts = &connectivity.Options{}
)

// createCmd represents the create command.
var createCmd = &cobra.Command{
	Use:   "create [machine name]",
	Short: "Create a new machine instance",
	Long: `Create a new machine instance on the specified cloud provider. 
The command requires a machine name and various flags to configure the instance properties 
such as provider, credentials, instance type, image, region and zone.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		createRequest.Name = args[0]

		providerName, err := cmd.Flags().GetString("provider")
		if err != nil {
			return err
		}

		credentialsName, err := cmd.Flags().GetString("credentials")
		if err != nil {
			return err
		}

		prov, err := provisioner.NewProvisioner(providerName, credentialsName)
		if err != nil {
			return fmt.Errorf("failed to create provisioner: %w", err)
		}

		err = prov.CreateMachine(cmd.Context(), createRequest, connectivityOpts)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	createCmd.Flags().StringP("credentials", "c", "", "The cloud provider credentials to use")
	_ = createCmd.MarkFlagRequired("credentials")
	createCmd.Flags().StringP("provider", "p", "", "The cloud provider")
	_ = createCmd.MarkFlagRequired("provider")

	// Machine specific flags
	createCmd.Flags().StringVarP(&createRequest.InstanceType, "instance-type", "t", "", "The instance type (or maybe named flavor, droplet, vm depending of the cloud provider)")
	createCmd.Flags().StringVarP(&createRequest.Image, "image", "i", "", "The image to use for the instance")
	createCmd.Flags().StringVarP(&createRequest.Region, "region", "r", "", "The cloud provider region where the instance will be spawned")
	createCmd.Flags().StringVarP(&createRequest.Zone, "zone", "z", "", "The zone in the region where the instance will be spawned")

	// Networking customization flags
	createCmd.Flags().BoolVar(&connectivityOpts.Public, "public", false, "Defines if the Ollama instance should be publicly exposed or not (not recommended), if set false you can use SSH tunnel or tailscale to connect to your Ollama instance.")
	createCmd.Flags().StringVar(&connectivityOpts.TailscaleAuthKey, "tailscale-auth-key", "", "The Tailscale authentication key to use for the instance")
}
