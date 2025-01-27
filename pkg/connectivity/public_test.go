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

package connectivity_test

import (
	"testing"

	"github.com/alexandrevilain/ollama-machine/pkg/cloudinit"
	"github.com/alexandrevilain/ollama-machine/pkg/connectivity"
	"github.com/alexandrevilain/ollama-machine/pkg/machine"
	"github.com/alexandrevilain/ollama-machine/pkg/provider"
	. "github.com/onsi/gomega"
)

func TestPublicProviderName(t *testing.T) {
	g := NewWithT(t)
	provider := &connectivity.PublicProvider{}
	g.Expect(provider.Name()).To(Equal("public"))
}

func TestPublicProviderInstallViaCloudInit(t *testing.T) {
	g := NewWithT(t)
	provider := &connectivity.PublicProvider{}
	cloudInitConfig := &cloudinit.Config{}

	provider.InstallViaCloudInit(cloudInitConfig)
	g.Expect(cloudInitConfig.RunCmd).To(HaveLen(1))
	g.Expect(cloudInitConfig.RunCmd[0]).To(Equal([]string{"sh", "-c", `echo "OLLAMA_HOST=0.0.0.0" > /home/ollama-machine/env`}))
}

func TestPublicProviderRetrieveOllamaHost(t *testing.T) {
	tests := map[string]struct {
		machine *machine.Machine
		want    string
	}{
		"retrieve host": {
			machine: &machine.Machine{
				Machine: &provider.Machine{
					IP: "1.2.3.4",
				},
			},
			want: "1.2.3.4",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			provider := &connectivity.PublicProvider{}
			host, err := provider.RetrieveOllamaHost(tt.machine)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(host).To(Equal(tt.want))
		})
	}
}
