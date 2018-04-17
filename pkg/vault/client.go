// Copyright (C) 2018 Nicolas Lamirault <nicolas.lamirault@gmail.com>

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vault

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/golang/glog"
	vaultapi "github.com/hashicorp/vault/api"

	pkgalan "github.com/nlamirault/alan/pkg/alan"
)

const (
	// DefaultAddr define the default Vault URL
	DefaultAddr = "http://127.0.0.1:8200"
)

// Client is the Client for the REST API of Vault
type Client struct {
	vault    *vaultapi.Client
	username string
	password string
}

// NewClient creates a client to manage Vault entities
func NewClient(vaultAddr string, username string, password string) (*Client, error) {
	glog.V(2).Infof("Setup using Vault server: %s %s", vaultAddr, username)
	transport := &http.Transport{}
	transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	vaultConfig := vaultapi.Config{
		Address: vaultAddr,
		HttpClient: &http.Client{
			Transport: transport,
		},
	}
	client, err := vaultapi.NewClient(&vaultConfig)
	if err != nil {
		return nil, err
	}
	return &Client{
		vault:    client,
		username: username,
		password: password,
	}, nil
}

// Login performs authentication with the Vault server
func (client *Client) Login() error {
	glog.V(2).Infof("Do authentication to the Vault using: %s", client.username)
	options := map[string]interface{}{
		"password": client.password,
	}

	// the login path
	path := fmt.Sprintf("auth/userpass/login/%s", client.username)
	secret, err := client.vault.Logical().Write(path, options)
	if err != nil {
		return err
	}
	client.vault.SetToken(secret.Auth.ClientToken)
	return nil
}

// Write create a new secret
func (client *Client) Write(key string, secret pkgalan.Secret) error {
	glog.V(2).Infof("Write secret: %s %s", key, secret)
	path := fmt.Sprintf("/secret/alan/%s", key)
	_, err := client.vault.Logical().Write(path,
		map[string]interface{}{
			"Title":    secret.Title,
			"URL":      secret.URL,
			"UserName": secret.Username,
			"Password": secret.Password,
		})
	if err != nil {
		return err
	}
	return nil
}

// Read retrieve a secret
func (client *Client) Read(key string) (map[string]interface{}, error) {
	glog.V(2).Infof("Read secret: %s ", key)
	path := fmt.Sprintf("/secret/alan/%s", key)
	secret, err := client.vault.Logical().Read(path)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("No secret for path %s", key)
	}
	return secret.Data, nil
}

// List retrieve some secrets
func (client *Client) List(key string) (map[string]interface{}, error) {
	glog.V(2).Infof("List secrets: %s ", key)
	path := fmt.Sprintf("/secret/alan/%s", key)
	secret, err := client.vault.Logical().List(path)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("No secrets for path %s", key)
	}
	return secret.Data, nil
}
