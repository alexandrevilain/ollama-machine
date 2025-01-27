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

	"github.com/alexandrevilain/ollama-machine/pkg/cloudcredentials"
	"github.com/alexandrevilain/ollama-machine/pkg/registry"
	"github.com/charmbracelet/log"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// credentialsCmd represents the credentials command.
var credentialsCmd = &cobra.Command{
	Use:   "credentials",
	Short: "Manage cloud provider credentials",
	Long:  `Manage cloud provider credentials used to create and manage machines.`,
}

var credentialsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all cloud credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := cloudcredentials.List()
		if err != nil {
			return err
		}

		table := uitable.New()
		table.MaxColWidth = 50

		table.AddRow("NAME", "PROVIDER")
		for _, cred := range result {
			table.AddRow(cred.Name, cred.Provider)
		}
		fmt.Println(table)

		return nil
	},
}

var credentialsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new cloud credential",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		credentialName := args[0]
		providerName, err := cmd.Flags().GetString("provider")
		if err != nil {
			return err
		}

		provider, found := registry.Providers[providerName]
		if !found {
			return fmt.Errorf("provider %q not found", providerName)
		}

		providerCredentials := provider.Credentials()

		err = providerCredentials.Complete()
		if err != nil {
			return err
		}

		err = providerCredentials.Validate()
		if err != nil {
			return err
		}

		err = cloudcredentials.Save(cloudcredentials.Key{Name: credentialName, Provider: providerName}, providerCredentials)
		if err != nil {
			if errors.Is(err, cloudcredentials.ErrAlreadyExists) {
				return fmt.Errorf("credential %q already exists, remove it first to update it", credentialName)
			}

			return err
		}

		log.Info("Cloud Credentials created")

		return nil
	},
}

var credentialsRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Remove a cloud credential",
	RunE: func(cmd *cobra.Command, args []string) error {
		credentialsName := args[0]

		providerName, err := cmd.Flags().GetString("provider")
		if err != nil {
			return err
		}

		err = cloudcredentials.Delete(cloudcredentials.Key{
			Name:     credentialsName,
			Provider: providerName,
		})
		if err != nil {
			return err
		}

		log.Info("Cloud Credentials removed")

		return nil
	},
}

func init() {
	credentialsCmd.PersistentFlags().StringP("provider", "p", "", "The provider to use for the credential")

	// Register flags for each provider.
	for providerName, provider := range registry.Providers {
		sub := pflag.NewFlagSet(providerName, pflag.ContinueOnError)
		sub.SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
			return pflag.NormalizedName(providerName + "-" + name)
		})
		provider.Credentials().RegisterFlags(sub)
		credentialsCreateCmd.Flags().AddFlagSet(sub)
	}

	_ = credentialsCreateCmd.MarkFlagRequired("provider")
	_ = credentialsRemoveCmd.MarkFlagRequired("provider")

	credentialsCmd.AddCommand(credentialsListCmd)
	credentialsCmd.AddCommand(credentialsRemoveCmd)
	credentialsCmd.AddCommand(credentialsCreateCmd)
}
