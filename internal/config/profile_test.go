// Copyright 2020 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build unit

package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCLIConfigHome(t *testing.T) {
	t.Run("with env set", func(t *testing.T) {
		expHome, err := os.UserConfigDir()
		expected := fmt.Sprintf("%s/mongocli", expHome)
		if err != nil {
			t.Fatalf("os.UserConfigDir() unexpected error: %v", err)
		}

		home, err := CLIConfigHome()
		if err != nil {
			t.Fatalf("MongoCLIConfigHome() unexpected error: %v", err)
		}
		if home != expected {
			t.Errorf("MongoCLIConfigHome() = %s; want '%s'", home, expected)
		}
	})
}

func TestOldMongoCLIConfigHome(t *testing.T) {
	t.Run("old home with XDG_CONFIG_HOME", func(t *testing.T) {
		const xdgHome = "my_config"
		t.Setenv("XDG_CONFIG_HOME", xdgHome)
		home, err := OldMongoCLIConfigHome()
		if err != nil {
			t.Fatalf("OldMongoCLIConfigHome() unexpected error: %v", err)
		}
		if home != xdgHome {
			t.Errorf("MongoCLIConfigHome() = %s; want '%s'", home, xdgHome)
		}
	})
	t.Run("old home without XDG_CONFIG_HOME", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")
		home, err := OldMongoCLIConfigHome()
		if err != nil {
			t.Fatalf("OldMongoCLIConfigHome() unexpected error: %v", err)
		}
		osHome, _ := os.UserHomeDir()
		if home != osHome+"/.config" {
			t.Errorf("OldMongoCLIConfigHome() = %s; want '%s/.config'", home, osHome)
		}
	})
}

func TestConfig_IsTrue(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{
			input: "true",
			want:  true,
		},
		{
			input: "True",
			want:  true,
		},
		{
			input: "t",
			want:  true,
		},
		{
			input: "T",
			want:  true,
		},
		{
			input: "TRUE",
			want:  true,
		},
		{
			input: "y",
			want:  true,
		},
		{
			input: "Y",
			want:  true,
		},
		{
			input: "yes",
			want:  true,
		},
		{
			input: "Yes",
			want:  true,
		},
		{
			input: "YES",
			want:  true,
		},
		{
			input: "1",
			want:  true,
		},
		{
			input: "false",
			want:  false,
		},
		{
			input: "f",
			want:  false,
		},
		{
			input: "unknown",
			want:  false,
		},
		{
			input: "0",
			want:  false,
		},
		{
			input: "",
			want:  false,
		},
	}
	for _, tt := range tests {
		if got := IsTrue(tt.input); got != tt.want {
			t.Errorf("IsTrue() get: %v, want %v", got, tt.want)
		}
	}
}

func TestConfig_SetName(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		require.NoError(t, SetName("default"))
		require.NoError(t, SetName("default-123"))
		require.NoError(t, SetName("default-test"))
	})

	t.Run("invalid", func(t *testing.T) {
		require.Error(t, SetName("d.efault"))
		require.Error(t, SetName("default.123"))
		require.Error(t, SetName("default.test"))
	})
}

func TestConfig_Rename(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		require.NoError(t, Rename("default"))
		require.NoError(t, Rename("default-123"))
		require.NoError(t, Rename("default-test"))
	})

	t.Run("invalid", func(t *testing.T) {
		require.Error(t, Rename("d.efault"))
		require.Error(t, Rename("default.123"))
		require.Error(t, Rename("default.test"))
	})
}
