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
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/flag"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/mocks"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/standby"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func failoverTestOpts(store standby.StateStore, out, errW *bytes.Buffer) *FailoverOpts {
	return &FailoverOpts{
		opts:       opts{store: store},
		OutputOpts: cli.OutputOpts{OutWriter: out, Template: failoverTemplate},
		confirm:    true,
		errW:       errW,
	}
}

func standbyStateDoc() *standby.RemoteDRState {
	return &standby.RemoteDRState{
		State:         standby.StateStandby,
		ClusterName:   "rs1",
		Version:       "1",
		LastModified:  "2026-08-19T12:25:23Z",
		SchemaVersion: "1",
	}
}

func TestFailover_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStateStore(ctrl)

	mockStore.EXPECT().FetchDRState(gomock.Any()).Return(standbyStateDoc(), `"etag-1"`, nil)
	mockStore.EXPECT().ConditionalPutDRState(gomock.Any(), gomock.Any(), `"etag-1"`).DoAndReturn(
		func(_ context.Context, s *standby.RemoteDRState, _ string) error {
			assert.Equal(t, standby.StatePromoteStandby, s.State)
			return nil
		})

	var out bytes.Buffer
	require.NoError(t, failoverTestOpts(mockStore, &out, &bytes.Buffer{}).Run(context.Background()))
	assert.Contains(t, out.String(), "Failover triggered")
	assert.Contains(t, out.String(), "rs1")
	assert.Contains(t, out.String(), standby.StatePromoteStandby)
}

func TestFailover_RunRejectsActive(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStateStore(ctrl)

	current := standbyStateDoc()
	current.State = standby.StateActive
	mockStore.EXPECT().FetchDRState(gomock.Any()).Return(current, `"etag-1"`, nil)
	// no write expected: any non-Standby state is refused

	err := failoverTestOpts(mockStore, &bytes.Buffer{}, &bytes.Buffer{}).Run(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "a failover can only be triggered from")
}

func TestFailover_RunWatchUntilActive(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStateStore(ctrl)

	active := standbyStateDoc()
	active.State = standby.StateActive
	active.PreviousState = standby.StatePromoteStandby
	active.Version = "3"

	gomock.InOrder(
		mockStore.EXPECT().FetchDRState(gomock.Any()).Return(standbyStateDoc(), `"etag-1"`, nil),
		mockStore.EXPECT().ConditionalPutDRState(gomock.Any(), gomock.Any(), `"etag-1"`).Return(nil),
		mockStore.EXPECT().FetchDRState(gomock.Any()).Return(active, `"etag-3"`, nil),
	)
	mockStore.EXPECT().FullPath().Return("s3://bucket/key").AnyTimes()

	prevInterval := watchPollInterval
	watchPollInterval = time.Millisecond
	t.Cleanup(func() { watchPollInterval = prevInterval })

	var out, errOut bytes.Buffer
	o := failoverTestOpts(mockStore, &out, &errOut)
	o.watch = true
	require.NoError(t, o.Run(context.Background()))
	assert.Contains(t, out.String(), "Failover completed")
	assert.Contains(t, errOut.String(), "Watching s3://bucket/key")
}

func TestFailoverBuilder(t *testing.T) {
	test.CmdValidator(
		t,
		FailoverBuilder(),
		0,
		[]string{
			flag.S3BucketName, flag.S3Key, flag.S3BucketEndpoint,
			flag.AWSRegion, flag.AWSAuthMode, flag.AWSAccessKey, flag.AWSSecretKey, flag.AWSRoleARN, flag.AWSProfile,
			flag.Force, flag.Watch, flag.Output,
		},
	)
}
