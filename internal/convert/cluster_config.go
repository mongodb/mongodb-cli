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

package convert

import (
	"fmt"
	"reflect"
	"strconv"

	"go.mongodb.org/ops-manager/opsmngr"
	"go.mongodb.org/ops-manager/search"
)

const (
	zero            = "0"
	one             = "1"
	file            = "file"
	fcvLessThanFour = "< 4.0"
	mongos          = "mongos"
)

// ClusterConfig configuration for a cluster
// This cluster can be used to patch an automation config.
type ClusterConfig struct {
	RSConfig `yaml:",inline"`
	MongoURI string           `yaml:"mongoURI,omitempty" json:"mongoURI,omitempty"`
	Shards   []*RSConfig      `yaml:"shards,omitempty" json:"shards,omitempty"`
	Config   *RSConfig        `yaml:"config,omitempty" json:"config,omitempty"`
	Mongos   []*ProcessConfig `yaml:"mongos,omitempty" json:"mongos,omitempty"`
}

// newReplicaSetCluster when config is a replicaset.
func newReplicaSetCluster(name string, s int) *ClusterConfig {
	rs := &ClusterConfig{}
	rs.Name = name
	rs.Processes = make([]*ProcessConfig, s)

	return rs
}

// newShardedCluster when config is a sharded cluster.
func newShardedCluster(s *opsmngr.ShardingConfig) *ClusterConfig {
	rs := &ClusterConfig{}
	rs.Name = s.Name
	rs.Shards = make([]*RSConfig, len(s.Shards))
	rs.Mongos = make([]*ProcessConfig, 0, 1)
	rs.Tags = s.Tags

	return rs
}

// PatchAutomationConfig adds the ClusterConfig to a opsmngr.AutomationConfig
// this method will modify the given AutomationConfig to add the new replica set or sharded cluster information.
func (c *ClusterConfig) PatchAutomationConfig(out *opsmngr.AutomationConfig) error {
	// A replica set should be just a list of processes
	if c.Processes != nil && c.Mongos == nil && c.Shards == nil && c.Config == nil {
		return c.patchReplicaSet(out)
	}
	// a sharded cluster will be a list of mongos (processes),
	// shards, each with a list of process (replica sets)
	// one (1) config server, with a list of process (replica set)
	if c.Processes == nil && c.Mongos != nil && c.Shards != nil && c.Config != nil {
		return c.patchSharding(out)
	}

	return ErrInvalidConfig
}

func (c *ClusterConfig) patchSharding(out *opsmngr.AutomationConfig) error {
	newCluster := newShardingConfig(c)
	// transform cli config to automation config
	for i, s := range c.Shards {
		s.Version = c.Version
		s.FeatureCompatibilityVersion = c.FeatureCompatibilityVersion
		if err := s.patchShard(out, c.Name); err != nil {
			return err
		}
		newCluster.Shards[i] = newShard(s)
	}
	c.Config.Version = c.Version
	c.Config.FeatureCompatibilityVersion = c.FeatureCompatibilityVersion
	if err := c.Config.patchConfigServer(out, c.Name); err != nil {
		return err
	}

	newProcesses := make([]*opsmngr.Process, len(c.Mongos))
	for i, pc := range c.Mongos {
		pc.ProcessType = mongos
		pc.setDefaults(&c.RSConfig)
		pc.setProcessName(out.Processes, c.Name, "mongos", strconv.Itoa(len(out.Processes)+i))
		newProcesses[i] = newMongosProcess(pc, c.Name)
	}
	// This value may not be present and is mandatory
	if out.Auth.DeploymentAuthMechanisms == nil {
		out.Auth.DeploymentAuthMechanisms = make([]string, 0)
	}
	patchProcesses(out, newCluster.Name, newProcesses)
	patchSharding(out, newCluster)
	return nil
}

func (c *ClusterConfig) addToMongoURI(p *opsmngr.Process) {
	if c.MongoURI == "" {
		c.MongoURI = fmt.Sprintf("mongodb://%s:%d", p.Hostname, p.Args26.NET.Port)
	} else {
		c.MongoURI = fmt.Sprintf("%s,%s:%d", c.MongoURI, p.Hostname, p.Args26.NET.Port)
	}
}

func newShard(rsConfig *RSConfig) *opsmngr.Shard {
	s := &opsmngr.Shard{
		ID: rsConfig.Name,
		RS: rsConfig.Name,
	}
	if s.Tags == nil {
		s.Tags = make([]string, 0)
	}
	return s
}

func newShardingConfig(c *ClusterConfig) *opsmngr.ShardingConfig {
	rs := &opsmngr.ShardingConfig{
		Name:                c.Name,
		Shards:              make([]*opsmngr.Shard, len(c.Shards)),
		ConfigServerReplica: c.Config.Name,
		Tags:                c.Tags,
		Draining:            make([]string, 0),
		Collections:         make([]*map[string]any, 0),
	}
	if rs.Tags == nil {
		rs.Tags = make([]*map[string]any, 0)
	}

	return rs
}

// patchProcesses replace replica set processes with new configuration
// this will disable all existing processes for the given replica set and remove the association
// Then try to patch then with the new config if one config exists for the same host:port.
func patchProcesses(out *opsmngr.AutomationConfig, newReplicaSetID string, newProcesses []*opsmngr.Process) {
	for i, oldProcess := range out.Processes {
		if oldProcess.Args26.Replication != nil && oldProcess.Args26.Replication.ReplSetName == newReplicaSetID {
			oldProcess.Disabled = true
			oldProcess.Args26.Replication = new(opsmngr.Replication)
		}
		oldName := oldProcess.Name
		pos, found := search.Processes(newProcesses, func(p *opsmngr.Process) bool {
			return p.Name == oldName
		})
		if found {
			keepSettings(oldProcess, newProcesses, pos)
			out.Processes[i] = newProcesses[pos]
			newProcesses = append(newProcesses[:pos], newProcesses[pos+1:]...)
		}
	}
	if len(newProcesses) > 0 {
		out.Processes = append(out.Processes, newProcesses...)
	}
}

// keepSettings preserves server-side fields that the CLI's describe/update
// roundtrip would otherwise silently drop. It walks the Process reflectively
// and, for every nil-able field (pointer, slice, map, interface), copies the
// old value onto the new process when the new value is nil. It recurses into
// struct and non-nil pointer-to-struct fields so nested optional blocks
// (e.g. Args26.NET.Compression, Args26.Storage.WiredTiger) are preserved too.
//
// Primitive scalars (string, int, bool, float) are deliberately NOT merged:
// the zero value is ambiguous with "user cleared this", and the CLI's
// ProcessConfig already roundtrips the scalar fields we surface, so they
// cannot be silently lost through this path.
//
// Trade-off: a user cannot unset a previously-set nil-able field by omitting
// it from the YAML — old will always win when new is nil. We accept this
// because the recurring bug is fields getting wiped, not preserved, and an
// explicit per-field allow-list has to be extended every time the SDK grows
// a new optional.
func keepSettings(oldProcess *opsmngr.Process, newProcesses []*opsmngr.Process, pos int) {
	mergeKeepNilable(
		reflect.ValueOf(newProcesses[pos]).Elem(),
		reflect.ValueOf(oldProcess).Elem(),
	)
}

// mergeKeepNilable copies src -> dst for every nil-able field on dst that is
// currently nil. dst and src must be the same struct type.
func mergeKeepNilable(dst, src reflect.Value) {
	if dst.Kind() != reflect.Struct || src.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < dst.NumField(); i++ {
		df, sf := dst.Field(i), src.Field(i)
		if df.CanSet() {
			mergeField(df, sf)
		}
	}
}

// mergeField dispatches a single struct field: copies src -> dst if dst is a
// nil nil-able, or recurses into struct / non-nil pointer-to-struct fields.
func mergeField(df, sf reflect.Value) {
	k := df.Kind()
	if k == reflect.Struct {
		mergeKeepNilable(df, sf)
		return
	}
	if k == reflect.Ptr && !df.IsNil() && !sf.IsNil() && df.Elem().Kind() == reflect.Struct {
		mergeKeepNilable(df.Elem(), sf.Elem())
		return
	}
	if isNilableKind(k) && df.IsNil() && !sf.IsNil() {
		df.Set(sf)
	}
}

func isNilableKind(k reflect.Kind) bool {
	return k == reflect.Ptr || k == reflect.Map || k == reflect.Slice || k == reflect.Interface
}

// patchReplicaSet patches the replica set if it exists, else adds it as a new replica set.
func patchReplicaSet(out *opsmngr.AutomationConfig, newReplicaSet *opsmngr.ReplicaSet) {
	pos, found := search.ReplicaSets(out.ReplicaSets, func(r *opsmngr.ReplicaSet) bool {
		return r.ID == newReplicaSet.ID
	})

	if !found {
		out.ReplicaSets = append(out.ReplicaSets, newReplicaSet)
		return
	}

	oldReplicaSet := out.ReplicaSets[pos]
	lastID := oldReplicaSet.Members[len(oldReplicaSet.Members)-1].ID
	for j, newMember := range newReplicaSet.Members {
		newHost := newMember.Host
		k, found := search.Members(oldReplicaSet.Members, func(m opsmngr.Member) bool {
			return m.Host == newHost
		})
		if found {
			newMember.ID = oldReplicaSet.Members[k].ID
		} else {
			lastID++
			newMember.ID = lastID
		}
		newReplicaSet.Members[j] = newMember
	}
	out.ReplicaSets[pos] = newReplicaSet
}

// patchSharding patches the shard if it exists, else adds it as a new shard.
func patchSharding(out *opsmngr.AutomationConfig, s *opsmngr.ShardingConfig) {
	pos, found := search.ShardingConfig(out.Sharding, func(r *opsmngr.ShardingConfig) bool {
		return r.Name == s.Name
	})
	if !found {
		out.Sharding = append(out.Sharding, s)
		return
	}

	out.Sharding[pos] = s
}
