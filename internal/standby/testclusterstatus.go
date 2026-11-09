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

//go:build unit || integration

package standby

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
)

// testClusterStatusFixture is the shared cluster_status_v1 BSON fixture, embedded
// at build time so callers can load it regardless of their working directory.
//
//go:embed testdata/cluster_status_v1.bson
var testClusterStatusFixture []byte

// LoadTestClusterStatusBytes returns the raw cluster_status_v1 BSON fixture bytes
// (estimatedDataLossSeconds 40, one shard rs0 with lagSeconds 30). Tests that seed
// the fixture into blob storage use this to keep the published bytes verbatim.
func LoadTestClusterStatusBytes(t testing.TB) []byte {
	t.Helper()
	return bytes.Clone(testClusterStatusFixture)
}

// LoadTestClusterStatus returns the shared cluster_status_v1 BSON fixture decoded
// into a ClusterStatus.
func LoadTestClusterStatus(t testing.TB) *ClusterStatus {
	t.Helper()
	status, err := UnmarshalClusterStatus(testClusterStatusFixture)
	require.NoError(t, err)
	return status
}
