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

package cloudinit_test

import (
	"testing"

	"github.com/alexandrevilain/ollama-machine/pkg/cloudinit"
	. "github.com/onsi/gomega"
)

func TestNewConfig(t *testing.T) {
	g := NewWithT(t)
	config := cloudinit.NewConfig()
	g.Expect(config).NotTo(BeNil())
	g.Expect(config.Hostname).To(BeEmpty())
	g.Expect(config.SSHAuthorizedKeys).To(BeEmpty())
	g.Expect(config.Users).To(BeEmpty())
	g.Expect(config.RunCmd).To(BeEmpty())
	g.Expect(config.WriteFiles).To(BeEmpty())
}

func TestAddSSHAuthorizedKeys(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		keys []string
		want int
	}{
		"single key": {
			keys: []string{"ssh-rsa AAAA..."},
			want: 1,
		},
		"multiple keys": {
			keys: []string{"ssh-rsa AAAA...", "ssh-rsa BBBB..."},
			want: 2,
		},
		"empty keys": {
			keys: []string{},
			want: 0,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			config := cloudinit.NewConfig()
			config.AddSSHAuthorizedKeys(tt.keys...)
			g.Expect(config.SSHAuthorizedKeys).To(HaveLen(tt.want))
			if tt.want > 0 {
				g.Expect(config.SSHAuthorizedKeys).To(Equal(tt.keys))
			}
		})
	}
}

func TestAddUser(t *testing.T) {
	tests := map[string]struct {
		user cloudinit.User
	}{
		"basic user": {
			user: cloudinit.User{
				Name:   "testuser",
				Groups: "sudo",
				Shell:  "/bin/bash",
			},
		},
		"user with SSH keys": {
			user: cloudinit.User{
				Name:              "testuser",
				SSHAuthorizedKeys: []string{"ssh-rsa AAAA..."},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			config := cloudinit.NewConfig()
			config.AddUser(tt.user)
			g.Expect(config.Users).To(HaveLen(1))
			g.Expect(config.Users[0]).To(Equal(tt.user))
		})
	}
}

func TestAddRunCmd(t *testing.T) {
	tests := map[string]struct {
		command []string
	}{
		"simple command": {
			command: []string{"echo hello"},
		},
		"complex command": {
			command: []string{"sh", "-c", "curl -fsSL https://example.com | bash"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			config := cloudinit.NewConfig()
			config.AddRunCmd(tt.command)
			g.Expect(config.RunCmd).To(HaveLen(1))
			g.Expect(config.RunCmd[0]).To(Equal(tt.command))
		})
	}
}

func TestAddFile(t *testing.T) {
	tests := map[string]struct {
		file cloudinit.File
	}{
		"basic file": {
			file: cloudinit.File{
				Path:    "/etc/test.conf",
				Content: "test content",
			},
		},
		"file with permissions": {
			file: cloudinit.File{
				Path:        "/etc/secure.conf",
				Content:     "secure content",
				Owner:       "root:root",
				Permissions: "0600",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			config := cloudinit.NewConfig()
			config.AddFile(tt.file)
			g.Expect(config.WriteFiles).To(HaveLen(1))
			g.Expect(config.WriteFiles[0]).To(Equal(tt.file))
		})
	}
}

func TestMarshal(t *testing.T) {
	g := NewWithT(t)
	config := cloudinit.NewConfig()
	config.Hostname = "testhost"
	config.AddSSHAuthorizedKeys("ssh-rsa AAAA...")
	config.AddUser(cloudinit.User{
		Name:   "testuser",
		Groups: "sudo",
	})
	config.AddRunCmd([]string{"echo hello"})
	config.AddFile(cloudinit.File{
		Path:    "/etc/test.conf",
		Content: "test content",
	})

	data, err := config.Marshal()
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(data).NotTo(BeNil())

	// Check for expected YAML content
	yamlStr := string(data)
	g.Expect(yamlStr).To(ContainSubstring("hostname: testhost"))
	g.Expect(yamlStr).To(ContainSubstring("ssh-rsa AAAA..."))
	g.Expect(yamlStr).To(ContainSubstring("name: testuser"))
	g.Expect(yamlStr).To(ContainSubstring("echo hello"))
	g.Expect(yamlStr).To(ContainSubstring("/etc/test.conf"))
}

func TestRender(t *testing.T) {
	g := NewWithT(t)
	config := cloudinit.NewConfig()
	config.Hostname = "testhost"

	data, err := config.Render()
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(data).NotTo(BeNil())

	// Check for #cloud-config header
	rendered := string(data)
	g.Expect(rendered).To(HavePrefix("#cloud-config\n"))
	g.Expect(rendered).To(ContainSubstring("hostname: testhost"))
}
