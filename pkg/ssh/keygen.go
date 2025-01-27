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

package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path"
	"runtime"

	gossh "golang.org/x/crypto/ssh"
)

const (
	rsaBitSize = 4096
)

// KeyPair represents a private & public keypair.
type KeyPair struct {
	PrivateKey []byte
	PublicKey  []byte
}

// KeyPairFiles holds the paths to a private & public key files.
type KeyPairFiles struct {
	PrivateKeyPath string `json:"publicKeyPath"`
	PublicKeyPath  string `json:"privateKeyPath"`
}

// NewKeyPair generates a new SSH keypair
// This will return a private & public key encoded as DER.
func NewKeyPair() (*KeyPair, error) {
	priv, err := rsa.GenerateKey(rand.Reader, rsaBitSize)
	if err != nil {
		return nil, fmt.Errorf("failed generating key: %w", err)
	}

	if err := priv.Validate(); err != nil {
		return nil, fmt.Errorf("failed validating key: %w", err)
	}

	privDer := x509.MarshalPKCS1PrivateKey(priv)

	pubSSH, err := gossh.NewPublicKey(&priv.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed creating public key: %w", err)
	}

	return &KeyPair{
		PrivateKey: privDer,
		PublicKey:  gossh.MarshalAuthorizedKey(pubSSH),
	}, nil
}

// WriteToFile writes keypair to files.
func (kp *KeyPair) WriteToFile(privateKeyPath string, publicKeyPath string) error {
	files := []struct {
		File  string
		Type  string
		Value []byte
	}{
		{
			File:  privateKeyPath,
			Value: pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Headers: nil, Bytes: kp.PrivateKey}),
		},
		{
			File:  publicKeyPath,
			Value: kp.PublicKey,
		},
	}

	for _, file := range files {
		f, err := os.Create(file.File)
		if err != nil {
			return fmt.Errorf("failed creating file: %w", err)
		}

		if _, err := f.Write(file.Value); err != nil {
			return fmt.Errorf("failed writing to file: %w", err)
		}

		// windows does not support chmod
		switch runtime.GOOS {
		case "darwin", "freebsd", "linux", "openbsd":
			if err := f.Chmod(0o600); err != nil { //nolint:mnd
				return err
			}
		}
	}

	return nil
}

// GenerateSSHKey generates SSH keypair based on a base path.
func GenerateSSHKey(basePath string, keyName string) (*KeyPair, *KeyPairFiles, error) {
	if _, err := os.Stat(basePath); err != nil {
		return nil, nil, err
	}

	keyPair, err := NewKeyPair()
	if err != nil {
		return nil, nil, fmt.Errorf("failed generating key pair: %w", err)
	}

	files := &KeyPairFiles{
		PrivateKeyPath: path.Join(basePath, keyName),
		PublicKeyPath:  path.Join(basePath, fmt.Sprintf("%s.pub", keyName)),
	}

	if err := keyPair.WriteToFile(files.PrivateKeyPath, files.PublicKeyPath); err != nil {
		return nil, nil, fmt.Errorf("failed writing keys to file(s): %w", err)
	}

	return keyPair, files, nil
}

// DeleteKeyPairFiles deletes the keypair files.
func DeleteKeyPairFiles(files *KeyPairFiles) error {
	if err := os.Remove(files.PrivateKeyPath); err != nil {
		return err
	}

	if err := os.Remove(files.PublicKeyPath); err != nil {
		return err
	}

	return nil
}
