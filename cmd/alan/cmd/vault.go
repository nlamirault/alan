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

	"github.com/golang/glog"
	"github.com/spf13/cobra"

	pkgalan "github.com/nlamirault/alan/pkg/alan"
	pkgcmd "github.com/nlamirault/alan/pkg/cmd"
	"github.com/nlamirault/alan/pkg/vault"
)

var (
	path string
)

type vaultCmd struct {
	out io.Writer
}

func newVaultCmd(out io.Writer) *cobra.Command {
	vaultCmd := &vaultCmd{
		out: out,
	}

	cmd := &cobra.Command{
		Use:   "vault",
		Short: "Manage Vault. See subcommands",
		RunE:  nil,
	}

	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get a secret under a path",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(path) == 0 {
				return fmt.Errorf("missing path")
			}
			vaultClient, err := vault.NewClient(vaultAddress, "alan", "turing")
			if err != nil {
				return err
			}
			return vaultCmd.get(vaultClient)
		},
	}
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List secrets under a path",
		RunE: func(cmd *cobra.Command, args []string) error {
			// if len(path) == 0 {
			// 	return fmt.Errorf("missing path")
			// }
			vaultClient, err := vault.NewClient(vaultAddress, "alan", "turing")
			if err != nil {
				return err
			}
			return vaultCmd.list(vaultClient)
		},
	}

	getCmd.PersistentFlags().StringVar(&path, "path", "", "Vault path")
	getCmd.PersistentFlags().StringVar(&vaultAddress, "vault", vault.DefaultAddr, "Vault address")
	listCmd.PersistentFlags().StringVar(&path, "path", "", "Vault path")
	listCmd.PersistentFlags().StringVar(&vaultAddress, "vault", vault.DefaultAddr, "Vault address")
	cmd.AddCommand(getCmd)
	cmd.AddCommand(listCmd)
	return cmd
}

func (cmd vaultCmd) get(vaultClient *vault.Client) error {
	glog.V(1).Infof("Get secret for path %s", path)
	if err := vaultClient.Login(); err != nil {
		return err
	}
	data, err := vaultClient.Read(path)
	if err != nil {
		return err
	}
	glog.V(1).Infof("Vault secret: %s", data)
	fmt.Printf("Username: %s\nPassword: %s\nURL: %s\n",
		pkgcmd.GreenOut(data[pkgalan.Username].(string)),
		pkgcmd.GreenOut(data[pkgalan.Password].(string)),
		pkgcmd.GreenOut(data[pkgalan.URL].(string)))
	return nil
}

func (cmd vaultCmd) list(vaultClient *vault.Client) error {
	glog.V(1).Infof("List secrets for path %s", path)
	if err := vaultClient.Login(); err != nil {
		return err
	}
	data, err := vaultClient.List(path)
	if err != nil {
		return err
	}
	glog.V(1).Infof("Vault secrets: %s", data)
	for _, key := range data["keys"].([]interface{}) {
		fmt.Printf("- %s\n", pkgcmd.GreenOut(key))
	}
	return nil
}
