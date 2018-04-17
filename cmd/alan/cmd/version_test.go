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
	"fmt"
	"strings"
	"testing"

	"github.com/nlamirault/alan/pkg/version"
)

func Test_VersionCommand(t *testing.T) {
	out := new(bytes.Buffer)
	cmd := newVersionCmd(out, "unit test version command")
	if err := cmd.Execute(); err != nil {
		t.Fatalf(err.Error())
	}
	text := out.String()
	if !strings.HasSuffix(text, fmt.Sprintf("v%s\n", version.Version)) {
		t.Fatalf("Invalid version: %s %s", text, version.Version)
	}
}
