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

package noop

import (
	"errors"

	"github.com/spf13/pflag"
)

type Credentials struct {
	Nothing string `json:"nothing"`
}

func (c *Credentials) Complete() error {
	return nil
}

func (c *Credentials) Validate() error {
	if c.Nothing == "" {
		return errors.New("nothing is required")
	}

	return nil
}

func (c *Credentials) RegisterFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.Nothing, "nothing", "", "nothing is nothing")
}
