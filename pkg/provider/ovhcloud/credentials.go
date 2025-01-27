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

package ovhcloud

import (
	"errors"

	"github.com/spf13/pflag"
)

type Credentials struct {
	Endpoint          string `json:"endpoint"`
	ApplicationKey    string `json:"applicationKey"`
	ApplicationSecret string `json:"applicationSecret"`
	ConsumerKey       string `json:"consumerKey"`
	ProjectID         string `json:"projectId"`
}

func (c *Credentials) Complete() error {
	return nil
}

func (c *Credentials) Validate() error {
	if c.Endpoint == "" {
		return errors.New("endpoint is required")
	}

	if c.ApplicationKey == "" {
		return errors.New("application key is required")
	}

	if c.ApplicationSecret == "" {
		return errors.New("application secret is required")
	}

	if c.ConsumerKey == "" {
		return errors.New("consumer key is required")
	}

	if c.ProjectID == "" {
		return errors.New("service name is required")
	}

	return nil
}

func (c *Credentials) RegisterFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.Endpoint, "endpoint", "ovh-eu", "OVHcloud API endpoint")
	fs.StringVar(&c.ApplicationKey, "application-key", "", "OVHcloud application key")
	fs.StringVar(&c.ApplicationSecret, "application-secret", "", "OVHcloud application secret")
	fs.StringVar(&c.ConsumerKey, "consumer-key", "", "OVHcloud consumer key")
	fs.StringVar(&c.ProjectID, "project-id", "", "OVHcloud cloud project ID (also named service name)")
}
