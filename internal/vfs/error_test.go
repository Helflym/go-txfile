// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
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

package vfs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/unix"
)

func TestErrorFmt(t *testing.T) {
	type testCase struct {
		err      *Error
		expected string
	}

	cases := map[string]testCase{
		"with op+kind only": testCase{
			err:      Err("pkg/op", ErrPermission, "", nil),
			expected: "pkg/op: permission denied",
		},
		"with path": testCase{
			err:      Err("pkg/op", ErrNotExist, "/path/to/file", nil),
			expected: "pkg/op: /path/to/file: file does not exist",
		},
		"nested os error": testCase{
			err: Err("pkg/op", ErrPermission, "path/to/file", &os.PathError{
				Op:   "stat",
				Path: "path/to/file",
				Err:  unix.EPERM,
			}),
			expected: "pkg/op: path/to/file: permission denied: stat path/to/file: operation not permitted",
		},
	}

	for name, test := range cases {
		test := test
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.err.Error())
		})
	}
}
