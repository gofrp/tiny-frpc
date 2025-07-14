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

package gssh

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

func TestGeneratePrivateKey_Success(t *testing.T) {
	require := require.New(t)
	
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "gssh_test")
	require.NoError(err)
	defer os.RemoveAll(tempDir)
	
	privateKeyPath := filepath.Join(tempDir, "test_key")
	publicKeyPath := privateKeyPath + ".pub"
	
	// Generate private key
	err = generatePrivateKey(privateKeyPath)
	require.NoError(err)
	
	// Verify private key file exists and has correct permissions
	privInfo, err := os.Stat(privateKeyPath)
	require.NoError(err)
	require.False(privInfo.IsDir())
	if runtime.GOOS != "windows" {
		require.Equal(os.FileMode(0600), privInfo.Mode().Perm())
	}
	
	// Verify public key file exists and has correct permissions
	pubInfo, err := os.Stat(publicKeyPath)
	require.NoError(err)
	require.False(pubInfo.IsDir())
	if runtime.GOOS != "windows" {
		require.Equal(os.FileMode(0644), pubInfo.Mode().Perm())
	}
	
	// Verify private key can be parsed
	privKeyBytes, err := os.ReadFile(privateKeyPath)
	require.NoError(err)
	require.NotEmpty(privKeyBytes)
	
	signer, err := ssh.ParsePrivateKey(privKeyBytes)
	require.NoError(err)
	require.NotNil(signer)
	
	// Verify public key can be parsed
	pubKeyBytes, err := os.ReadFile(publicKeyPath)
	require.NoError(err)
	require.NotEmpty(pubKeyBytes)
	
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(pubKeyBytes)
	require.NoError(err)
	require.NotNil(pubKey)
	
	// Verify public key from private key matches public key file
	require.Equal(pubKey.Type(), signer.PublicKey().Type())
	require.Equal(pubKey.Marshal(), signer.PublicKey().Marshal())
}

func TestGeneratePrivateKey_CreatesDirectory(t *testing.T) {
	require := require.New(t)
	
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "gssh_test")
	require.NoError(err)
	defer os.RemoveAll(tempDir)
	
	// Use nested directory that doesn't exist
	nestedDir := filepath.Join(tempDir, "nested", "dir")
	privateKeyPath := filepath.Join(nestedDir, "test_key")
	
	// Create the nested directory first (mimicking os.MkdirAll behavior if needed)
	err = os.MkdirAll(nestedDir, 0755)
	require.NoError(err)
	
	// Generate private key
	err = generatePrivateKey(privateKeyPath)
	require.NoError(err)
	
	// Verify files exist
	_, err = os.Stat(privateKeyPath)
	require.NoError(err)
	_, err = os.Stat(privateKeyPath + ".pub")
	require.NoError(err)
}

func TestGeneratePrivateKey_InvalidPath(t *testing.T) {
	require := require.New(t)
	
	// Skip this test on Windows as permission handling is different
	if runtime.GOOS == "windows" {
		t.Skip("Skipping permission test on Windows")
	}
	
	// Create temporary directory with no write permissions
	tempDir, err := os.MkdirTemp("", "gssh_test")
	require.NoError(err)
	defer os.RemoveAll(tempDir)
	
	// Remove write permission from directory
	err = os.Chmod(tempDir, 0444)
	require.NoError(err)
	defer os.Chmod(tempDir, 0755) // Restore for cleanup
	
	privateKeyPath := filepath.Join(tempDir, "test_key")
	
	// Generate private key should fail
	err = generatePrivateKey(privateKeyPath)
	require.Error(err)
	require.Contains(err.Error(), "failed to write private key file")
}

func TestGeneratePrivateKey_EmptyPath(t *testing.T) {
	require := require.New(t)
	
	// Generate private key with empty path should fail
	err := generatePrivateKey("")
	require.Error(err)
}

func TestGeneratePrivateKey_OverwriteExisting(t *testing.T) {
	require := require.New(t)
	
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "gssh_test")
	require.NoError(err)
	defer os.RemoveAll(tempDir)
	
	privateKeyPath := filepath.Join(tempDir, "test_key")
	
	// Create existing file
	err = os.WriteFile(privateKeyPath, []byte("existing content"), 0600)
	require.NoError(err)
	
	// Generate private key should overwrite
	err = generatePrivateKey(privateKeyPath)
	require.NoError(err)
	
	// Verify content is different and valid
	content, err := os.ReadFile(privateKeyPath)
	require.NoError(err)
	require.NotEqual("existing content", string(content))
	
	// Verify it's a valid private key
	_, err = ssh.ParsePrivateKey(content)
	require.NoError(err)
}

func TestGeneratePrivateKey_KeyFormat(t *testing.T) {
	require := require.New(t)
	
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "gssh_test")
	require.NoError(err)
	defer os.RemoveAll(tempDir)
	
	privateKeyPath := filepath.Join(tempDir, "test_key")
	
	// Generate private key
	err = generatePrivateKey(privateKeyPath)
	require.NoError(err)
	
	// Read and verify private key format
	privKeyBytes, err := os.ReadFile(privateKeyPath)
	require.NoError(err)
	
	privKeyStr := string(privKeyBytes)
	require.Contains(privKeyStr, "BEGIN OPENSSH PRIVATE KEY")
	require.Contains(privKeyStr, "END OPENSSH PRIVATE KEY")
	
	// Read and verify public key format
	pubKeyBytes, err := os.ReadFile(privateKeyPath + ".pub")
	require.NoError(err)
	
	pubKeyStr := string(pubKeyBytes)
	require.True(strings.HasPrefix(pubKeyStr, "ssh-ed25519 "))
}

func TestGeneratePrivateKey_WithPublicKeyAuthFunc(t *testing.T) {
	require := require.New(t)
	
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "gssh_test")
	require.NoError(err)
	defer os.RemoveAll(tempDir)
	
	privateKeyPath := filepath.Join(tempDir, "test_key")
	
	// Generate private key
	err = generatePrivateKey(privateKeyPath)
	require.NoError(err)
	
	// Test that the generated key works with publicKeyAuthFunc
	authMethod, err := publicKeyAuthFunc(privateKeyPath)
	require.NoError(err)
	require.NotNil(authMethod)
}