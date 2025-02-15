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

	"github.com/alexandrevilain/ollama-machine/pkg/machine"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command.
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all machines",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		machines, err := machine.List()
		if err != nil {
			return err
		}

		table := uitable.New()
		table.MaxColWidth = 50

		table.AddRow("NAME", "STATE", "PROVIDER", "REGION", "IP", "OLLAMA HOST", "OLLAMA PORT")
		for _, machine := range machines {
			table.AddRow(machine.Name, machine.State, machine.ProviderName, machine.Region, machine.IP, machine.OllamaConfig.Host, machine.OllamaConfig.Port)
		}

		fmt.Println(table)

		return nil
	},
}
