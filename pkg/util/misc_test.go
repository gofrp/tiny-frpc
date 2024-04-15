package util

import (
	"reflect"
	"strings"
	"testing"
)

func TestJSONEncode(t *testing.T) {
	tests := []struct {
		name      string
		args      interface{}
		want      string
		wantError bool
	}{
		{
			name: "normal struct",
			args: struct {
				Name string
				Age  int
			}{"John", 30},
			want:      `{"Name":"John","Age":30}`,
			wantError: false,
		},
		{
			name:      "invalid data",
			args:      make(chan int), // channels are not serializable to JSON
			want:      "args: {}, error: json: unsupported type: chan int",
			wantError: true,
		},
		{
			name: "normal map",
			args: map[string]interface{}{
				"hello": "world",
				"age":   25,
			},
			want:      `{"age":25,"hello":"world"}`,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := JSONEncode(tt.args)
			if (got != tt.want) != tt.wantError {
				t.Errorf("JSONEncode() = %v, want %v", got, tt.want)
			}

			if tt.wantError && !strings.Contains(got, "error") {
				t.Errorf("Expected an error in JSONEncode() for input %v, but got none", tt.args)
			}
		})
	}
}

func TestEmptyOrInt(t *testing.T) {
	tests := []struct {
		value    int
		fallback int
		want     int
	}{
		{0, 10, 10},
		{5, 10, 5},
	}

	for _, tt := range tests {
		if got := EmptyOr(tt.value, tt.fallback); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("EmptyOr<int>() with input '%v': got %v, want %v", tt.value, got, tt.want)
		}
	}
}

func TestEmptyOrString(t *testing.T) {
	tests := []struct {
		value    string
		fallback string
		want     string
	}{
		{"", "fallback", "fallback"},
		{"value", "fallback", "value"},
	}

	for _, tt := range tests {
		if got := EmptyOr(tt.value, tt.fallback); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("EmptyOr<string>() with input '%v': got %v, want %v", tt.value, got, tt.want)
		}
	}
}

func TestEmptyOrBool(t *testing.T) {
	tests := []struct {
		value    bool
		fallback bool
		want     bool
	}{
		{false, true, true},
		{true, false, true},
	}

	for _, tt := range tests {
		if got := EmptyOr(tt.value, tt.fallback); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("EmptyOr<bool>() with input '%v': got %v, want %v", tt.value, got, tt.want)
		}
	}
}

func TestIsJSONBuffer(t *testing.T) {
	tests := []struct {
		buf  []byte
		want bool
	}{
		{[]byte(`{"key": "value"}`), true},
		{[]byte(` `), false},
		{[]byte(`not JSON`), false},
		{[]byte(`   {"key": "value"}`), true},
		{[]byte(`[]`), false},
		{nil, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.buf), func(t *testing.T) {
			if got := IsJSONBuffer(tt.buf); got != tt.want {
				t.Errorf("IsJSONBuffer() = %v, want %v for buf %q", got, tt.want, tt.buf)
			}
		})
	}
}

func TestTernary(t *testing.T) {
	tests := []struct {
		name       string
		condition  bool
		ifOutput   interface{}
		elseOutput interface{}
		want       interface{}
	}{
		{
			name:       "Condition true int",
			condition:  true,
			ifOutput:   1,
			elseOutput: 2,
			want:       1,
		},
		{
			name:       "Condition false int",
			condition:  false,
			ifOutput:   1,
			elseOutput: 2,
			want:       2,
		},
		{
			name:       "Condition true string",
			condition:  true,
			ifOutput:   "TrueValue",
			elseOutput: "FalseValue",
			want:       "TrueValue",
		},
		{
			name:       "Condition false string",
			condition:  false,
			ifOutput:   "TrueValue",
			elseOutput: "FalseValue",
			want:       "FalseValue",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Ternary(tc.condition, tc.ifOutput, tc.elseOutput)
			if got != tc.want {
				t.Errorf("Ternary() = %v, want %v", got, tc.want)
			}
		})
	}
}
