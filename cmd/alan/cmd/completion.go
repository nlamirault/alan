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

	"github.com/spf13/cobra"

	pkgcmd "github.com/nlamirault/alan/pkg/cmd"
)

var (
	completionShells = map[string]func(out io.Writer, cmd *cobra.Command) error{
		"bash": runCompletionBash,
		"zsh":  runCompletionZsh,
	}
)

func newCompletionCmd(out io.Writer, example string) *cobra.Command {
	shells := []string{}
	for s := range completionShells {
		shells = append(shells, s)
	}

	cmd := &cobra.Command{
		Use:     "completion SHELL",
		Short:   "Output shell completion code for the specified shell (bash or zsh)",
		Example: example,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := RunCompletion(out, cmd, args); err != nil {
				fmt.Fprintln(out, pkgcmd.RedOut(err))
			}
			return nil
		},
		ValidArgs: shells,
	}

	return cmd
}

func RunCompletion(out io.Writer, cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Shell not specified")
	}
	if len(args) > 1 {
		return fmt.Errorf("Too many arguments. Expected only the shell type: %s", args)
	}
	run, found := completionShells[args[0]]
	if !found {
		return fmt.Errorf("Unsupported shell type %q", args[0])
	}

	return run(out, cmd.Parent())
}

func runCompletionBash(out io.Writer, cmd *cobra.Command) error {
	return cmd.GenBashCompletion(out)
}

func runCompletionZsh(out io.Writer, cmd *cobra.Command) error {
	return fmt.Errorf("Zsh is currently Unsupported")
}
