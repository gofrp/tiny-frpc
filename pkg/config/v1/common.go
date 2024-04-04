// Copyright 2023 The frp Authors
//
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

package v1

import (
	"fmt"
	"sync"
)

// TODO(fatedier): Due to the current implementation issue of the go json library, the UnmarshalJSON method
// of a custom struct cannot access the DisallowUnknownFields parameter of the parent decoder.
// Here, a global variable is temporarily used to control whether unknown fields are allowed.
// Once the v2 version is implemented by the community, we can switch to a standardized approach.
//
// https://github.com/golang/go/issues/41144
// https://github.com/golang/go/discussions/63397
var (
	DisallowUnknownFields   = false
	DisallowUnknownFieldsMu sync.Mutex
)

type AuthMethod string

const (
	AuthMethodToken AuthMethod = "token"
)

type HeaderOperations struct {
	Set map[string]string `json:"set,omitempty"`
}

// ValidatePort checks that the network port is in range
func ValidatePort(port int, fieldPath string) error {
	if 0 <= port && port <= 65535 {
		return nil
	}
	return fmt.Errorf("%s: port number %d must be in the range 0..65535", fieldPath, port)
}
