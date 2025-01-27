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

package cloudcredentials

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	service  = "ollama-machine"
	indexKey = "index"
)

// ErrAlreadyExists is returned when trying to save a key that already exists.
var ErrAlreadyExists = errors.New("key already exists")

// Get retrieves the value associated with the given key and unmarshals it into dest.
func Get(key Key, dest any) error {
	result, err := keyring.Get(service, keyName(key))
	if err != nil {
		return fmt.Errorf("can't get key content: %w", err)
	}

	err = json.Unmarshal([]byte(result), dest)
	if err != nil {
		return fmt.Errorf("can't unmarshal key content: %w", err)
	}

	return nil
}

// Save saves the given value associated with the given key.
func Save(key Key, value any) error {
	if exists(key) {
		return ErrAlreadyExists
	}

	content, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("can't marshal value: %w", err)
	}

	err = keyring.Set(service, keyName(key), string(content))
	if err != nil {
		return fmt.Errorf("can't save key: %w", err)
	}

	return addKeyToIndex(key)
}

// Delete deletes the value associated with the given key.
func Delete(key Key) error {
	err := keyring.Delete(service, keyName(key))
	if err != nil {
		return fmt.Errorf("can't delete key: %w", err)
	}

	return removeKeyFromIndex(key)
}

// List lists all credentials in the store.
func List() (StoreIndex, error) {
	index, err := keyring.Get(service, indexKey)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return StoreIndex{}, nil
		}

		return nil, fmt.Errorf("can't get index content: %w", err)
	}

	result := StoreIndex{}
	err = json.Unmarshal([]byte(index), &result)

	return result, err
}

// exists checks if the given key exists in the index.
func exists(key Key) bool {
	index, err := List()
	if err != nil {
		return false
	}

	for _, e := range index {
		if e.Name == key.Name && e.Provider == key.Provider {
			return true
		}
	}

	return false
}

// addKeyToIndex adds the given key to the index.
func addKeyToIndex(key Key) error {
	index, err := List()
	if err != nil {
		return fmt.Errorf("can't get index: %w", err)
	}

	index = append(index, key)
	content, err := json.Marshal(index)
	if err != nil {
		return fmt.Errorf("can't marshal index: %w", err)
	}

	err = keyring.Set(service, indexKey, string(content))
	if err != nil {
		return fmt.Errorf("can't save index: %w", err)
	}

	return nil
}

// removeKeyFromIndex removes the given key from the index.
func removeKeyFromIndex(key Key) error {
	index, err := List()
	if err != nil {
		return fmt.Errorf("can't get index: %w", err)
	}

	for i, e := range index {
		if e.Name == key.Name && e.Provider == key.Provider {
			index = append(index[:i], index[i+1:]...)

			break
		}
	}

	content, err := json.Marshal(index)
	if err != nil {
		return fmt.Errorf("can't marshal index: %w", err)
	}
	err = keyring.Set(service, indexKey, string(content))
	if err != nil {
		return fmt.Errorf("can't save index: %w", err)
	}

	return nil
}

// keyName returns the key name for the given entry.
func keyName(entry Key) string {
	return fmt.Sprintf("%s-%s", entry.Name, entry.Provider)
}
