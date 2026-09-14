// Copyright 2026 MongoDB Inc
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

package keyproviders

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func writeSelfSignedCA(t *testing.T, path string, withPrivateKey bool) {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "test-ca"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	require.NoError(t, err)

	contents := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	if withPrivateKey {
		contents = append(contents, pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		})...)
	}
	require.NoError(t, os.WriteFile(path, contents, 0600))
}

func TestValidateServerCA(t *testing.T) {
	t.Run("CA file with only a certificate block is valid", func(t *testing.T) {
		caFile := filepath.Join(t.TempDir(), "ca.pem")
		writeSelfSignedCA(t, caFile, false)

		ki := &KMIPKeyIdentifier{ServerCAFileName: caFile}

		require.NoError(t, ki.validateServerCA())
	})

	t.Run("CA file with certificate and private key blocks is valid", func(t *testing.T) {
		caFile := filepath.Join(t.TempDir(), "ca.pem")
		writeSelfSignedCA(t, caFile, true)

		ki := &KMIPKeyIdentifier{ServerCAFileName: caFile}

		require.NoError(t, ki.validateServerCA())
	})

	t.Run("CA file without a certificate block is invalid", func(t *testing.T) {
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		require.NoError(t, err)

		caFile := filepath.Join(t.TempDir(), "ca.pem")
		keyOnly := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		})
		require.NoError(t, os.WriteFile(caFile, keyOnly, 0600))

		ki := &KMIPKeyIdentifier{ServerCAFileName: caFile}

		require.ErrorContains(t, ki.validateServerCA(), "server CA file does not contain a certificate block")
	})
}
