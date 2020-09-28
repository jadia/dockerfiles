// Copyright 2018 The ksonnet authors
//
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package clicmd

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/ksonnet/ksonnet/pkg/actions"
)

func Test_applyCmd(t *testing.T) {
	cases := []cmdTestCase{
		{
			name:   "with no options",
			args:   []string{"apply", "default"},
			action: actionApply,
			expected: map[string]interface{}{
				actions.OptionApp:            mock.AnythingOfType("*app.App"),
				actions.OptionEnvName:        "default",
				actions.OptionGcTag:          "",
				actions.OptionSkipGc:         false,
				actions.OptionComponentNames: make([]string, 0),
				actions.OptionCreate:         true,
				actions.OptionDryRun:         false,
				actions.OptionClientConfig:   mock.AnythingOfType("*client.Config"),
			},
		},
		{
			name:  "invalid jsonnet flag",
			args:  []string{"apply", "default", "--ext-str", "foo"},
			isErr: true,
		},
	}

	runTestCmd(t, cases)
}
