// Copyright 2023 MongoDB Inc
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

// This code was autogenerated at 2023-04-25T17:59:17+01:00. Note: Manual updates are allowed, but may be overwritten.

package atlas

import (
	atlasv2 "go.mongodb.org/atlas-sdk/v20231115004/admin"
)

//go:generate mockgen -destination=../../mocks/atlas/mock_data_lake_pipelines_runs.go -package=atlas github.com/mongodb/mongodb-atlas-cli/internal/store/atlas PipelineRunsLister,PipelineRunsDescriber

type PipelineRunsLister interface {
	PipelineRuns(string, string) (*atlasv2.PaginatedPipelineRun, error)
}

type PipelineRunsDescriber interface {
	PipelineRun(string, string, string) (*atlasv2.IngestionPipelineRun, error)
}

// PipelineRuns encapsulates the logic to manage different cloud providers.
func (s *Store) PipelineRuns(projectID, pipelineName string) (*atlasv2.PaginatedPipelineRun, error) {
	result, _, err := s.clientv2.DataLakePipelinesApi.ListPipelineRuns(s.ctx, projectID, pipelineName).Execute()
	return result, err
}

// PipelineRun encapsulates the logic to manage different cloud providers.
func (s *Store) PipelineRun(projectID, pipelineName, id string) (*atlasv2.IngestionPipelineRun, error) {
	result, _, err := s.clientv2.DataLakePipelinesApi.GetPipelineRun(s.ctx, projectID, pipelineName, id).Execute()
	return result, err
}
