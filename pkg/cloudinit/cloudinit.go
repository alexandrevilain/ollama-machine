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

// package cloudinit provides functionality to create and marshal cloud-init configurations
package cloudinit

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Config represents the main cloud-init configuration structure.
type Config struct {
	Hostname          string     `yaml:"hostname,omitempty"`
	SSHAuthorizedKeys []string   `yaml:"ssh_authorized_keys,omitempty"` //nolint:tagliatelle
	Users             []User     `yaml:"users,omitempty"`
	RunCmd            [][]string `yaml:"runcmd,omitempty"`
	Bootcmd           []string   `yaml:"bootcmd,omitempty"`
	WriteFiles        []File     `yaml:"write_files,omitempty"` //nolint:tagliatelle
}

// User represents a user configuration.
type User struct {
	Name              string   `yaml:"name"`
	Groups            string   `yaml:"groups,omitempty"`
	Shell             string   `yaml:"shell,omitempty"`
	Sudo              string   `yaml:"sudo,omitempty"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"` //nolint:tagliatelle
	PasswdHash        string   `yaml:"passwd,omitempty"`
}

// File represents a file to be written.
type File struct {
	Path        string `yaml:"path"`
	Content     string `yaml:"content"`
	Owner       string `yaml:"owner,omitempty"`
	Permissions string `yaml:"permissions,omitempty"`
	Encoding    string `yaml:"encoding,omitempty"`
}

// NewConfig creates a new cloud-init configuration.
func NewConfig() *Config {
	return &Config{}
}

// AddSSHAuthorizedKeys adds SSH authorized keys to the configuration.
func (c *Config) AddSSHAuthorizedKeys(keys ...string) {
	c.SSHAuthorizedKeys = append(c.SSHAuthorizedKeys, keys...)
}

// AddUser adds a new user to the configuration.
func (c *Config) AddUser(user User) {
	c.Users = append(c.Users, user)
}

// AddRunCmd adds a command to be run.
func (c *Config) AddRunCmd(cmd []string) {
	c.RunCmd = append(c.RunCmd, cmd)
}

// AddFile adds a file to be written.
func (c *Config) AddFile(file File) {
	c.WriteFiles = append(c.WriteFiles, file)
}

// Marshal returns the YAML representation of the configuration.
func (c *Config) Marshal() ([]byte, error) {
	return yaml.Marshal(c)
}

// Render returns the bytes representation of the configuration.
func (c *Config) Render() ([]byte, error) {
	data, err := c.Marshal()
	if err != nil {
		return nil, fmt.Errorf("error marshaling config: %w", err)
	}

	b := bytes.NewBufferString("#cloud-config\n")
	_, err = b.Write(data)

	return b.Bytes(), err
}
