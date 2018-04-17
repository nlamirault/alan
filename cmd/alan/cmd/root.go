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
	goflag "flag"
	"fmt"
	"io"
	"os"

	_ "github.com/golang/glog" // init glog to get its flags
	"github.com/spf13/cobra"

	pkgcmd "github.com/nlamirault/alan/pkg/cmd"
)

var (
	cliName           = "alan"
	helpMessage       = "alan - Bridge between Vault and password managers"
	completionExample = `
               # Load the alan completion code for bash into the current shell
               source <(alan completion bash)

               # Write bash completion code to a file and source if from .bash_profile
               alan completion bash > ~/.alan/completion.bash.inc
               printf "\n# Picous shell completion\nsource '$HOME/.alan/completion.bash.inc'\n" >> $HOME/.bash_profile
               source $HOME/.bash_profile

               # Load the alan completion code for zsh[1] into the current shell
			   source <(alan completion zsh)`

	vaultAddress string
)

func newApplicationCommand(out io.Writer) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          cliName,
		Long:         `Bridge between Vault and password managers`,
		SilenceUsage: true,
	}
	rootCmd.AddCommand(
		newVersionCmd(out, helpMessage),
		newCompletionCmd(out, completionExample),
		newKeepassXCCmd(out),
		newVaultCmd(out),
	)
	cobra.EnablePrefixMatching = true

	// add glog flags
	rootCmd.PersistentFlags().AddGoFlagSet(goflag.CommandLine)
	// https://github.com/kubernetes/dns/pull/27/files
	goflag.CommandLine.Parse([]string{})

	return rootCmd
}

func Execute() {
	cmd := newApplicationCommand(os.Stdout)
	if err := cmd.Execute(); err != nil {
		fmt.Println(pkgcmd.RedOut(err))
		os.Exit(1)
	}
}
