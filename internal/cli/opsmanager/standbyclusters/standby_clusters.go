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

package standbyclusters

import (
	"cmp"
	"context"
	"errors"
	"fmt"

	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/cli"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/config"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/flag"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/standby"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/usage"
	"github.com/spf13/cobra"
)

const (
	// Resource name of this command subtree.
	Use              = "standby-clusters"
	defaultAWSRegion = "us-east-1"
)

var errMissingStateFileKey = errors.New("no DR state file configured. Set it with 'mongocli ops-manager standby-clusters configure' or pass --s3Key")

func Builder() *cobra.Command {
	cmd := &cobra.Command{
		Use:     Use,
		Aliases: []string{"standbyClusters", "sc"},
		Short:   "Manage standby cluster disaster recovery without Ops Manager.",
		Long: `These commands read and write the standby cluster disaster recovery (DR)
state file directly in Amazon S3. They don't use the Ops Manager API, so they
remain available when Ops Manager is unreachable.

Run 'mongocli ops-manager standby-clusters configure' first to set up access
to the DR state file of your standby cluster.`,
	}

	cmd.AddCommand(
		ConfigureBuilder(),
		DescribeBuilder(),
		FailoverBuilder(),
	)

	return cmd
}

type opts struct {
	cli.GlobalOpts

	bucket     string
	key        string
	region     string
	authMode   string
	accessKey  string
	secretKey  string
	roleARN    string
	awsProfile string
	endpoint   string

	store standby.StateStore
}

func (o *opts) rawCredentials() standby.Credentials {
	return standby.Credentials{
		AuthMode:        standby.AuthMode(cmp.Or(o.authMode, config.StandbyAWSAuthMode())),
		RoleArn:         cmp.Or(o.roleARN, config.StandbyAWSRoleARN()),
		AccessKeyID:     cmp.Or(o.accessKey, config.StandbyAWSAccessKeyID()),
		SecretAccessKey: cmp.Or(o.secretKey, config.StandbyAWSSecretAccessKey()),
		BucketName:      cmp.Or(o.bucket, config.StandbyS3Bucket()),
		Region:          cmp.Or(o.region, config.StandbyAWSRegion()),
		AWSProfile:      cmp.Or(o.awsProfile, config.StandbyAWSProfile()),
		Endpoint:        cmp.Or(o.endpoint, config.StandbyS3Endpoint()),
	}
}

// credentials resolves the effective AWS credentials from flags and profile.
func (o *opts) credentials() standby.Credentials {
	c := o.rawCredentials()
	// Same inference as the agent: role ARN set -> assumeRole, access key set
	// -> staticCredentials, otherwise the default AWS credentials chain.
	if c.AuthMode == "" {
		switch {
		case c.RoleArn != "":
			c.AuthMode = standby.AuthModeAssumeRole
		case c.AccessKeyID != "":
			c.AuthMode = standby.AuthModeStaticCredentials
		default:
			c.AuthMode = standby.AuthModeCredentialsChain
		}
	}
	if c.Region == "" {
		c.Region = defaultAWSRegion
	}
	c.UsePathStyle = c.Endpoint != ""
	return c
}

func (o *opts) stateFileKey() string {
	return cmp.Or(o.key, config.StandbyS3Key())
}

// initStore builds the S3 state store from the resolved credentials. Used as
// a PreRunE step by the commands that talk to the DR state file.
func (o *opts) initStore(ctx context.Context) func() error {
	return func() error {
		key := o.stateFileKey()
		if key == "" {
			return errMissingStateFileKey
		}
		store, err := standby.NewS3StateStore(ctx, o.credentials(), key)
		if err != nil {
			return fmt.Errorf("%w (fix with 'mongocli ops-manager standby-clusters configure' or the --s3*/--aws* flags)", err)
		}
		o.store = store
		return nil
	}
}

func addSharedFlags(cmd *cobra.Command, o *opts) {
	cmd.Flags().StringVar(&o.bucket, flag.S3BucketName, "", usage.StandbyS3BucketName)
	cmd.Flags().StringVar(&o.key, flag.S3Key, "", usage.StandbyS3Key)
	cmd.Flags().StringVar(&o.region, flag.AWSRegion, "", usage.StandbyAWSRegion)
	cmd.Flags().StringVar(&o.authMode, flag.AWSAuthMode, "", usage.StandbyAWSAuthMode)
	cmd.Flags().StringVar(&o.accessKey, flag.AWSAccessKey, "", usage.AWSAccessKey)
	cmd.Flags().StringVar(&o.secretKey, flag.AWSSecretKey, "", usage.AWSSecretKey)
	cmd.Flags().StringVar(&o.roleARN, flag.AWSRoleARN, "", usage.StandbyAWSRoleARN)
	cmd.Flags().StringVar(&o.awsProfile, flag.AWSProfile, "", usage.StandbyAWSProfile)
	cmd.Flags().StringVar(&o.endpoint, flag.S3BucketEndpoint, "", usage.StandbyS3Endpoint)
}
