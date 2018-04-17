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

	"github.com/nlamirault/alan/pkg/version"
)

type versionCmd struct {
	out io.Writer
}

func newVersionCmd(out io.Writer, help string) *cobra.Command {
	versionCmd := &versionCmd{
		out: out,
	}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number.",
		Long:  `All software has versions. This is alan's.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return versionCmd.printVersion(help)
		},
	}
	return cmd
}

func (cmd versionCmd) printVersion(help string) error {
	fmt.Fprintf(cmd.out, "%s. v%s\n", help, version.Version)
	return nil
}
