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

package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/golang/glog"
	"github.com/spf13/cobra"

	pkgalan "github.com/nlamirault/alan/pkg/alan"
	pkgcmd "github.com/nlamirault/alan/pkg/cmd"
	"github.com/nlamirault/alan/pkg/keepassxc"
	"github.com/nlamirault/alan/pkg/vault"
)

var (
	database string
)

type keepassxcCmd struct {
	out io.Writer
}

func newKeepassXCCmd(out io.Writer) *cobra.Command {
	keepassxcCmd := &keepassxcCmd{
		out: out,
	}

	cmd := &cobra.Command{
		Use:   "keepassxc",
		Short: "Manage KeepassXC database. See subcommands",
		RunE:  nil,
	}

	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Show a KeepassXC database",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(database) == 0 {
				return fmt.Errorf("missing database name")
			}
			return keepassxcCmd.showDB()
		},
	}
	importCmd := &cobra.Command{
		Use:   "import",
		Short: "Import a KeepassXC database into a Vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(database) == 0 {
				return fmt.Errorf("missing database name")
			}
			vaultClient, err := vault.NewClient(vaultAddress, "alan", "turing")
			if err != nil {
				return err
			}
			return keepassxcCmd.importDB(vaultClient)
		},
	}
	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "Export Vault entries to a KeepassXC database",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(database) == 0 {
				return fmt.Errorf("missing database name")
			}
			vaultClient, err := vault.NewClient(vaultAddress, "alan", "turing")
			if err != nil {
				return err
			}
			return keepassxcCmd.exportDB(vaultClient)
		},
	}

	showCmd.PersistentFlags().StringVar(&database, "database", "", "Database filename")
	importCmd.PersistentFlags().StringVar(&database, "database", "", "Database filename")
	importCmd.PersistentFlags().StringVar(&vaultAddress, "vault", vault.DefaultAddr, "Vault address")
	exportCmd.PersistentFlags().StringVar(&database, "database", "", "Database filename")
	exportCmd.PersistentFlags().StringVar(&vaultAddress, "vault", vault.DefaultAddr, "Vault address")
	cmd.AddCommand(showCmd)
	cmd.AddCommand(importCmd)
	cmd.AddCommand(exportCmd)
	return cmd
}

func (cmd keepassxcCmd) importDB(vaultClient *vault.Client) error {
	glog.V(1).Infof("Import database: %s", database)
	keepassClient, err := keepassxc.NewClient(database)
	if err != nil {
		return err
	}
	if err := keepassClient.Open(); err != nil {
		return err
	}
	if err := vaultClient.Login(); err != nil {
		return err
	}
	secrets, err := keepassClient.Load()
	if err != nil {
		return err
	}
	//glog.V(2).Infof("Secrets for Vault: %s", secrets)
	for name, group := range secrets {
		glog.V(2).Infof("Manage Vault group: %s", name)
		for _, secret := range group {
			glog.V(2).Infof("Manage Vault secret: %s", secret)
			if len(secret.Title) == 0 {
				fmt.Printf(pkgcmd.YellowOut(fmt.Sprintf("No title for secret: %s %s\n", secret.Username, secret.URL)))
			} else {
				path := fmt.Sprintf("%s/%s", name, secret.Title)
				fmt.Printf(pkgcmd.GreenOut(fmt.Sprintf("Add secret: %s\n", path)))
				vaultClient.Write(path, secret)
			}
		}
	}
	return keepassClient.Close()
}

func (cmd keepassxcCmd) showDB() error {
	glog.V(1).Infof("Show database: %s", database)
	keepassClient, err := keepassxc.NewClient(database)
	if err != nil {
		return err
	}
	if err := keepassClient.Open(); err != nil {
		return err
	}
	secrets, err := keepassClient.Load()
	if err != nil {
		return err
	}
	for name, group := range secrets {
		fmt.Println(pkgcmd.GreenOut(name))
		for _, secret := range group {
			if len(secret.Title) == 0 {
				fmt.Printf("%s %s %s\n", pkgcmd.RedOut(">>>"), pkgcmd.YellowOut(secret.Username), pkgcmd.YellowOut(secret.URL))
			} else {
				fmt.Printf("%s: %s %s\n", pkgcmd.BlueOut(secret.Title), pkgcmd.BlueOut(secret.Username), pkgcmd.BlueOut(secret.URL))
			}
		}
	}

	return keepassClient.Close()
}

func (cmd keepassxcCmd) exportDB(vaultClient *vault.Client) error {
	glog.V(1).Infof("Export database: %s", database)
	if err := vaultClient.Login(); err != nil {
		return err
	}

	keepassClient, err := keepassxc.NewClient(database)
	if err != nil {
		return err
	}

	data, err := vaultClient.List(path)
	if err != nil {
		return err
	}
	glog.V(1).Infof("Vault secrets: %s", data)
	secrets := map[string][]*pkgalan.Secret{}
	for _, key := range data["keys"].([]interface{}) {
		path := key.(string)
		keySecrets, err := extractKeyEntries(keepassClient, vaultClient, path)
		if err != nil {
			return err
		}
		secrets[path] = keySecrets
	}

	if err := keepassClient.Create(secrets); err != nil {
		return err
	}
	return keepassClient.Save()
}

func extractKeyEntries(keepassClient *keepassxc.Client, vaultClient *vault.Client, path string) ([]*pkgalan.Secret, error) {
	glog.V(2).Infof("Analyse secrets for path: %s", path)
	secrets := []*pkgalan.Secret{}
	glog.V(2).Infof("Analyse Vault group: %s", path)
	fmt.Printf(pkgcmd.GreenOut(fmt.Sprintf("Vault group: %s\n", path)))
	data, err := vaultClient.List(path)
	if err != nil {
		return nil, err
	}
	for _, key := range data["keys"].([]interface{}) {
		newPath := fmt.Sprintf("%s%s", path, key.(string))
		if strings.HasSuffix(key.(string), "/") {
			glog.V(2).Infof("Analyse sub entry: %s", newPath)
			keySecrets, err := extractKeyEntries(keepassClient, vaultClient, newPath)
			if err != nil {
				return nil, err
			}
			secrets = append(secrets, keySecrets...)
		} else {
			glog.V(2).Infof("Secret for: %s", newPath)
			fmt.Printf(pkgcmd.BlueOut(fmt.Sprintf("Vault entry: %s\n", newPath)))
			keyData, err := vaultClient.Read(newPath)
			if err != nil {
				return nil, err
			}
			glog.V(1).Infof("Vault secret: %s", keyData)
			if keyData["Title"] != nil {
				secrets = append(secrets, &pkgalan.Secret{
					Title:    keyData[pkgalan.Title].(string),
					Username: keyData[pkgalan.Username].(string),
					Password: keyData[pkgalan.Password].(string),
					URL:      keyData[pkgalan.URL].(string),
				})
			}
		}
	}

	return secrets, nil
}
