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
	"os"
	"text/template"

	"github.com/alexandrevilain/ollama-machine/pkg/machine"
	"github.com/docker/machine/libmachine/shell"
	"github.com/spf13/cobra"
)

const (
	envTmpl = `{{ .Prefix }}OLLAMA_HOST{{ .Delimiter }}{{ .OllamaHost }}{{ .Suffix }}`
)

type ShellConfig struct {
	Prefix     string
	Delimiter  string
	Suffix     string
	OllamaHost string
}

// envCmd represents the env command.
var envCmd = &cobra.Command{
	Use:   "env [machine name]",
	Short: "",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		machineName := args[0]
		shell, err := shell.Detect()
		if err != nil {
			return err
		}

		m, err := machine.GetByName(machineName)
		if err != nil {
			return err
		}

		shellCfg := ShellConfig{
			OllamaHost: m.OllamaConfig.Address(),
		}

		switch shell {
		case "fish":
			shellCfg.Prefix = "set -gx "
			shellCfg.Suffix = "\";\n"
			shellCfg.Delimiter = " \""
		case "powershell":
			shellCfg.Prefix = "$Env:"
			shellCfg.Suffix = "\"\n"
			shellCfg.Delimiter = " = \""
		case "cmd":
			shellCfg.Prefix = "SET "
			shellCfg.Suffix = "\n"
			shellCfg.Delimiter = "="
		case "tcsh":
			shellCfg.Prefix = "setenv "
			shellCfg.Suffix = "\";\n"
			shellCfg.Delimiter = " \""
		case "emacs":
			shellCfg.Prefix = "(setenv \""
			shellCfg.Suffix = "\")\n"
			shellCfg.Delimiter = "\" \""
		default:
			shellCfg.Prefix = "export "
			shellCfg.Suffix = "\"\n"
			shellCfg.Delimiter = "=\""
		}

		return executeTemplateStdout(&shellCfg)
	},
}

func executeTemplateStdout(shellCfg *ShellConfig) error {
	t := template.New("envConfig")
	tmpl, err := t.Parse(envTmpl)
	if err != nil {
		return err
	}

	return tmpl.Execute(os.Stdout, shellCfg)
}
