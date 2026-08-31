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

package standby_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/standby"

	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/mocks"
)

func standbyDoc() *standby.RemoteDRState {
	return &standby.RemoteDRState{
		State:         standby.StateStandby,
		ClusterName:   "rs1",
		Version:       "1",
		LastModified:  "2026-08-19T12:25:23Z",
		SchemaVersion: "1",
	}
}

func TestTriggerFailover(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockStateStore(ctrl)

	current := standbyDoc()
	store.EXPECT().FetchDRState(gomock.Any()).Return(current, `"etag-1"`, nil)
	store.EXPECT().ConditionalPutDRState(gomock.Any(), gomock.Any(), `"etag-1"`).DoAndReturn(
		func(_ context.Context, s *standby.RemoteDRState, _ string) error {
			assert.Equal(t, standby.StatePromoteStandby, s.State)
			assert.Equal(t, standby.StateStandby, s.PreviousState)
			assert.Equal(t, "2", s.Version)
			return nil
		})

	written, err := standby.TriggerFailover(context.Background(), store)
	require.NoError(t, err)
	assert.Equal(t, standby.StatePromoteStandby, written.State)
}

func TestTriggerFailover_AlreadyPromoting(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockStateStore(ctrl)

	current := standbyDoc()
	current.State = standby.StatePromoteStandby
	current.PreviousState = standby.StateStandby
	current.Version = "2"
	store.EXPECT().FetchDRState(gomock.Any()).Return(current, `"etag-1"`, nil)
	// No write expected: any non-Standby state is refused.

	_, err := standby.TriggerFailover(context.Background(), store)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "a failover can only be triggered from")
}

func TestTriggerFailover_MissingStateFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockStateStore(ctrl)

	store.EXPECT().FetchDRState(gomock.Any()).Return(nil, "", nil)
	store.EXPECT().FullPath().Return("s3://b/k").AnyTimes()

	_, err := standby.TriggerFailover(context.Background(), store)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no DR state file found")
}

func TestTriggerFailover_RejectsActive(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockStateStore(ctrl)

	current := standbyDoc()
	current.State = "Active"
	store.EXPECT().FetchDRState(gomock.Any()).Return(current, `"etag-1"`, nil)

	_, err := standby.TriggerFailover(context.Background(), store)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "a failover can only be triggered from")
}

func TestTriggerFailover_RetriesOnceOnConflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockStateStore(ctrl)

	first := standbyDoc()
	second := standbyDoc()
	second.Version = "2" // a concurrent writer bumped it

	gomock.InOrder(
		store.EXPECT().FetchDRState(gomock.Any()).Return(first, `"etag-1"`, nil),
		store.EXPECT().ConditionalPutDRState(gomock.Any(), gomock.Any(), `"etag-1"`).Return(standby.ErrPreconditionFailed),
		store.EXPECT().FetchDRState(gomock.Any()).Return(second, `"etag-2"`, nil),
		store.EXPECT().ConditionalPutDRState(gomock.Any(), gomock.Any(), `"etag-2"`).DoAndReturn(
			func(_ context.Context, s *standby.RemoteDRState, _ string) error {
				assert.Equal(t, "3", s.Version, "retry must recompute from the re-read state")
				return nil
			}),
	)

	written, err := standby.TriggerFailover(context.Background(), store)
	require.NoError(t, err)
	assert.Equal(t, "3", written.Version)
}

func TestTriggerFailover_DoubleConflictFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockStateStore(ctrl)

	store.EXPECT().FetchDRState(gomock.Any()).DoAndReturn(
		func(context.Context) (*standby.RemoteDRState, string, error) {
			return standbyDoc(), `"e"`, nil // fresh copy per call
		}).Times(2)
	store.EXPECT().ConditionalPutDRState(gomock.Any(), gomock.Any(), `"e"`).Return(standby.ErrPreconditionFailed).Times(2)
	store.EXPECT().FullPath().Return("s3://b/k").AnyTimes()

	_, err := standby.TriggerFailover(context.Background(), store)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "changed while writing")
}

func TestTriggerFailover_PutErrorPropagates(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockStateStore(ctrl)

	store.EXPECT().FetchDRState(gomock.Any()).Return(standbyDoc(), `"e"`, nil)
	store.EXPECT().ConditionalPutDRState(gomock.Any(), gomock.Any(), `"e"`).Return(errors.New("access denied"))

	_, err := standby.TriggerFailover(context.Background(), store)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "access denied")
}
