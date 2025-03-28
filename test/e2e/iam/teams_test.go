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

//go:build e2e || (iam && !om60 && !atlas)

package iam_test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/mongodb/mongodb-cli/mongocli/v2/test/e2e"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/ops-manager/opsmngr"
)

func TestTeams(t *testing.T) {
	cliPath, err := e2e.Bin()
	require.NoError(t, err)

	n, err := e2e.RandInt(1000)
	require.NoError(t, err)

	teamName := fmt.Sprintf("teams%v", n)
	var teamID string

	t.Run("Create", func(t *testing.T) {
		username, _, err := OrgNUser(0)
		require.NoError(t, err)

		cmd := exec.Command(cliPath,
			iamEntity,
			teamsEntity,
			"create",
			teamName,
			"--username",
			username,
			"-o=json")
		cmd.Env = os.Environ()
		resp, err := cmd.CombinedOutput()

		a := assert.New(t)
		require.NoError(t, err, string(resp))

		var team opsmngr.Team
		require.NoError(t, json.Unmarshal(resp, &team))
		a.Equal(teamName, team.Name)
		teamID = team.ID
	})
	require.NotEmpty(t, teamID)

	t.Run("Describe By ID", func(t *testing.T) {
		cmd := exec.Command(cliPath,
			iamEntity,
			teamsEntity,
			"describe",
			"--id",
			teamID,
			"-o=json")
		cmd.Env = os.Environ()
		resp, err := cmd.CombinedOutput()

		a := assert.New(t)
		require.NoError(t, err, string(resp))

		var team opsmngr.Team
		require.NoError(t, json.Unmarshal(resp, &team))
		a.Equal(teamID, team.ID)
	})

	t.Run("Describe By Name", func(t *testing.T) {
		cmd := exec.Command(cliPath,
			iamEntity,
			teamsEntity,
			"describe",
			"--name",
			teamName,
			"-o=json")
		cmd.Env = os.Environ()
		resp, err := cmd.CombinedOutput()

		a := assert.New(t)
		require.NoError(t, err, string(resp))

		var team opsmngr.Team
		require.NoError(t, json.Unmarshal(resp, &team))
		a.Equal(teamName, team.Name)
	})

	t.Run("List", func(t *testing.T) {
		cmd := exec.Command(cliPath,
			iamEntity,
			teamsEntity,
			"ls",
			"-o=json")
		cmd.Env = os.Environ()
		resp, err := cmd.CombinedOutput()

		a := assert.New(t)
		require.NoError(t, err, string(resp))

		var teams []opsmngr.Team
		require.NoError(t, json.Unmarshal(resp, &teams))
		a.NotEmpty(t, teams)
	})

	t.Run("Delete", func(t *testing.T) {
		cmd := exec.Command(cliPath,
			iamEntity,
			teamsEntity,
			"delete",
			teamID,
			"--force")
		cmd.Env = os.Environ()
		resp, err := cmd.CombinedOutput()
		a := assert.New(t)
		require.NoError(t, err, string(resp))
		expected := fmt.Sprintf("Team '%s' deleted\n", teamID)
		a.Equal(expected, string(resp))
	})
}
