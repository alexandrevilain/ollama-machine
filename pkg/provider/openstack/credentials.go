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

package openstack

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

type Credentials struct {
	IdentityEndpoint   string `json:"identityEndpoint"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	TenantID           string `json:"tenantId"`
	TenantName         string `json:"tenantName"`
	DomainName         string `json:"domainName"`
	Region             string `json:"region"`
	IdentityAPIVersion int    `json:"identityApiVersion"`

	passwordFromStdin bool `json:"-"`
}

func (c *Credentials) Complete() error {
	if c.passwordFromStdin {
		passwordFromStdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		password := strings.TrimSuffix(string(passwordFromStdin), "\n")
		password = strings.TrimSuffix(password, "\r")

		c.Password = password
	}

	return nil
}

func (c *Credentials) Validate() error {
	if c.IdentityEndpoint == "" {
		return errors.New("identity-endpoint is required")
	}

	if c.Username == "" {
		return errors.New("username is required")
	}

	if c.Password == "" {
		return errors.New("password is required")
	}

	if c.TenantID == "" && c.TenantName == "" {
		return errors.New("tenant-id or tenant-name is required")
	}

	if c.DomainName == "" {
		return errors.New("domain-name is required")
	}

	if c.Region == "" {
		return errors.New("region is required")
	}

	if c.IdentityAPIVersion != 2 && c.IdentityAPIVersion != 3 {
		return errors.New("identity-api-version must be 2 or 3")
	}

	return nil
}

func (c *Credentials) RegisterFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.IdentityEndpoint, "identity-endpoint", "", "OpenStack identity endpoint")
	fs.StringVar(&c.Username, "username", "", "OpenStack username")
	fs.StringVar(&c.Password, "password", "", "OpenStack password")
	fs.BoolVar(&c.passwordFromStdin, "password-from-stdin", false, "Read OpenStack password from stdin")
	fs.StringVar(&c.TenantID, "tenant-id", "", "OpenStack tenant ID")
	fs.StringVar(&c.TenantName, "tenant-name", "", "OpenStack tenant name")
	fs.StringVar(&c.DomainName, "domain-name", "Default", "OpenStack user domain name")
	fs.StringVar(&c.Region, "region", "", "OpenStack region")
	fs.IntVar(&c.IdentityAPIVersion, "identity-api-version", 3, "OpenStack identity API version") //nolint:mnd
}
