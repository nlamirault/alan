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
	"bytes"
	"testing"
)

func executeCompletionCommand(t *testing.T, args []string) string {
	out := new(bytes.Buffer)
	cmd := newCompletionCmd(out, "unit test completion command")
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Fatalf(err.Error())
	}
	return out.String()
}

func Test_CompletionCmdWithoutShell(t *testing.T) {
	text := executeCompletionCommand(t, []string{""})
	if len(text) == 0 {
		t.Fatalf("Invalid completion: %s", text)
	}
}

// func Test_CompletionCmdWithBash(t *testing.T) {
// 	text := executeCompletionCommand(t, []string{"bash"})
// 	if len(text) == 0 {
// 		t.Fatalf("Invalid completion: %s", text)
// 	}
// }
