// Copyright 2024 gofrp (https://github.com/gofrp)
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

package log

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		name   string
		level  string
		format string
		values []interface{}
	}{
		{
			name:   "INFO",
			level:  LevelInfo,
			format: "This is an %s message",
			values: []interface{}{"info"},
		},
		{
			name:   "WARN",
			level:  LevelWarn,
			format: "This is an %s message",
			values: []interface{}{"warning"},
		},
		{
			name:   "ERROR",
			level:  LevelError,
			format: "This is an %s message",
			values: []interface{}{"error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLogger(&buf)

			switch tt.level {
			case LevelInfo:
				logger.Info(tt.format, tt.values...)
			case LevelWarn:
				logger.Warn(tt.format, tt.values...)
			case LevelError:
				logger.Error(tt.format, tt.values...)
			}

			output := buf.String()
			if !strings.Contains(output, tt.level) || !strings.Contains(output, "This is an") {
				t.Errorf("log output should contain the level '%s' and message 'This is an', got: %s", tt.level, output)
			}

			// Optionally check if the output contains a timestamp
			if !strings.Contains(output, time.Now().Format(time.RFC3339)[:10]) {
				t.Errorf("log output should contain a correctly formatted timestamp, got: %s", output)
			}
		})
	}
}
