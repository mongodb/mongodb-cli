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

package standbyclusters

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/flag"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/mocks"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/standby"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func describeTestOpts(store standby.StateStore, out *bytes.Buffer) *DescribeOpts {
	return &DescribeOpts{
		opts:       opts{store: store},
		OutputOpts: cli.OutputOpts{OutWriter: out, Template: describeTemplate},
	}
}

func TestDescribe_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStateStore(ctrl)

	state := &standby.RemoteDRState{
		State:         standby.StateStandby,
		ClusterName:   "rs1",
		Version:       "1",
		LastModified:  "2026-08-19T12:25:23Z",
		SchemaVersion: "1",
	}
	mockStore.EXPECT().FetchDRState(gomock.Any()).Return(state, `"etag"`, nil).Times(1)

	var out bytes.Buffer
	require.NoError(t, describeTestOpts(mockStore, &out).Run(context.Background()))
	assert.Contains(t, out.String(), standby.StateStandby)
	assert.Contains(t, out.String(), "rs1")
}

func TestDescribe_RunNoStateFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStateStore(ctrl)

	mockStore.EXPECT().FetchDRState(gomock.Any()).Return(nil, "", nil).Times(1)
	mockStore.EXPECT().FullPath().Return("s3://bucket/key").AnyTimes()

	err := describeTestOpts(mockStore, &bytes.Buffer{}).Run(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no DR state file found")
}

func TestDescribe_RunFetchError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStateStore(ctrl)

	mockStore.EXPECT().FetchDRState(gomock.Any()).Return(nil, "", errors.New("boom")).Times(1)

	err := describeTestOpts(mockStore, &bytes.Buffer{}).Run(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "boom")
}

func TestDescribeBuilder(t *testing.T) {
	test.CmdValidator(
		t,
		DescribeBuilder(),
		0,
		[]string{
			flag.S3BucketName, flag.S3Key, flag.S3BucketEndpoint,
			flag.AWSRegion, flag.AWSAuthMode, flag.AWSAccessKey, flag.AWSSecretKey, flag.AWSRoleARN, flag.AWSProfile,
			flag.Output,
		},
	)
}

func TestBuilder(t *testing.T) {
	test.CmdValidator(
		t,
		Builder(),
		3,
		[]string{},
	)
}
