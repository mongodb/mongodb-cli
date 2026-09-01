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
	"context"
	"fmt"
	"io"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/config"
	"github.com/mongodb/mongodb-cli/mongocli/v2/internal/standby"
	"github.com/spf13/cobra"
)

type ConfigureOpts struct {
	opts

	out io.Writer
}

func (o *ConfigureOpts) Run(ctx context.Context) error {
	if o.out == nil {
		return fmt.Errorf("no output writer configured")
	}
	_, _ = fmt.Fprintf(o.out, `You are configuring access to the standby cluster disaster recovery (DR)
state file in Amazon S3. These settings are stored in your mongocli profile
and never touch the Ops Manager API.

`)

	if err := o.ask(); err != nil {
		return err
	}

	creds := o.credentials()
	if err := creds.Validate(); err != nil {
		return err
	}

	// Smoke test: the configuration is only useful if the file is readable.
	store, err := standby.NewS3StateStore(ctx, creds, o.stateFileKey())
	if err != nil {
		return err
	}
	state, _, err := store.FetchDRState(ctx)
	if err != nil {
		return fmt.Errorf("could not read the DR state file at %s: %w", store.FullPath(), err)
	}
	if state == nil {
		return fmt.Errorf("no DR state file found at %s; verify the bucket and key", store.FullPath())
	}
	_, _ = fmt.Fprintf(o.out, "Verified access to %s (current state: %s).\n", store.FullPath(), state.State)

	o.save()

	if err := config.Save(); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(o.out, `
Your standby cluster settings are saved in profile '%s'.
You can use [%s config set] to change these settings at a later time.
`, config.Name(), config.MongoCLI)
	return nil
}

func (o *ConfigureOpts) save() {
	config.SetStandbyS3Bucket(o.credentials().BucketName)
	config.SetStandbyS3Key(o.stateFileKey())
	config.SetStandbyAWSRegion(o.credentials().Region)
	config.SetStandbyAWSAuthMode(string(o.credentials().AuthMode))
	config.SetStandbyAWSAccessKeyID(o.credentials().AccessKeyID)
	config.SetStandbyAWSSecretAccessKey(o.credentials().SecretAccessKey)
	config.SetStandbyAWSRoleARN(o.credentials().RoleArn)
	config.SetStandbyAWSProfile(o.credentials().AWSProfile)
	config.SetStandbyS3Endpoint(o.credentials().Endpoint)
}

// ask prompts for every value not already provided via flags. Profile
// settings become the prompt defaults, so Enter keeps the current value.
func (o *ConfigureOpts) ask() error {
	c := o.rawCredentials()

	if err := o.askLocation(&c); err != nil {
		return err
	}

	if o.authMode == "" {
		if err := askAuthMode(&c); err != nil {
			return err
		}
	}

	if err := o.askModeCredentials(&c); err != nil {
		return err
	}

	// Write the prompted values back so credentials() picks them up.
	o.bucket, o.region, o.endpoint = c.BucketName, c.Region, c.Endpoint
	o.authMode = string(c.AuthMode)
	o.accessKey, o.secretKey = c.AccessKeyID, c.SecretAccessKey
	o.roleARN, o.awsProfile = c.RoleArn, c.AWSProfile
	return nil
}

// askLocation prompts for the state file location: bucket, key, region, endpoint.
func (o *ConfigureOpts) askLocation(c *standby.Credentials) error {
	if err := askInput(&c.BucketName, "S3 bucket name holding the DR state file (name only, no s3:// prefix or URL):", "", o.bucket != ""); err != nil {
		return err
	}
	key := o.stateFileKey()
	if err := askInput(&key, "DR state file object key within the bucket (e.g. <clusterPrefix>/dr_status_<componentId>.json):", "", o.key != ""); err != nil {
		return err
	}
	o.key = key

	regionDefault := c.Region
	if regionDefault == "" {
		regionDefault = defaultAWSRegion
	}
	if err := askInput(&c.Region, "AWS region of the bucket:", regionDefault, o.region != ""); err != nil {
		return err
	}
	return askInput(&c.Endpoint, "Custom S3 endpoint URL (leave empty for AWS S3; set for MinIO and other S3-compatible stores):", c.Endpoint, o.endpoint != "")
}

// askModeCredentials prompts for the fields required by the chosen auth mode.
func (o *ConfigureOpts) askModeCredentials(c *standby.Credentials) error {
	switch c.AuthMode {
	case standby.AuthModeStaticCredentials:
		if err := askInput(&c.AccessKeyID, "AWS access key ID:", "", o.accessKey != ""); err != nil {
			return err
		}
		if o.secretKey == "" {
			return survey.AskOne(&survey.Password{Message: "AWS secret access key:"}, &c.SecretAccessKey)
		}
	case standby.AuthModeAssumeRole:
		if err := askInput(&c.RoleArn, "AWS IAM role ARN to assume:", "", o.roleARN != ""); err != nil {
			return err
		}
		return askAWSProfile(c, o.awsProfile != "")
	case standby.AuthModeCredentialsChain:
		return askAWSProfile(c, o.awsProfile != "")
	}
	return nil
}

func askAWSProfile(c *standby.Credentials, skip bool) error {
	if skip {
		return nil
	}
	return survey.AskOne(&survey.Input{
		Message: "AWS profile to use (leave empty for the default chain, use an AWS SSO profile name after 'aws sso login'):",
		Default: c.AWSProfile,
	}, &c.AWSProfile)
}

func askAuthMode(c *standby.Credentials) error {
	modes := make([]string, 0, len(standby.ValidAuthModes()))
	for _, m := range standby.ValidAuthModes() {
		modes = append(modes, string(m))
	}
	answer := string(c.AuthMode)
	if answer == "" {
		answer = string(standby.AuthModeCredentialsChain)
	}
	if err := survey.AskOne(&survey.Select{
		Message: "AWS authentication mode:",
		Options: modes,
		Default: answer,
		Help:    "staticCredentials: access key + secret key. assumeRole: IAM role ARN over ambient credentials. credentialsChain: AWS SDK default chain — environment variables, ~/.aws, instance profile, or AWS SSO after 'aws sso login'.",
	}, &answer); err != nil {
		return err
	}
	c.AuthMode = standby.AuthMode(answer)
	return nil
}

// askInput prompts with an optional default. Skipped when the value was
// passed as a flag (scripted/non-interactive use).
func askInput(target *string, message, defaultValue string, skip bool) error {
	if skip {
		return nil
	}
	return survey.AskOne(&survey.Input{Message: message, Default: defaultValue}, target)
}

// mongocli ops-manager standby-clusters configure.
func ConfigureBuilder() *cobra.Command {
	o := &ConfigureOpts{}
	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Configure access to the disaster recovery state file of your standby cluster.",
		Long: `Guides you through setting up access to the standby cluster disaster
recovery (DR) state file in Amazon S3: bucket, object key, region, and AWS
credentials. Supports static AWS access keys, IAM role assumption, and the
default AWS credentials chain (environment variables, ~/.aws, instance
profiles, and AWS SSO after 'aws sso login').

The settings are stored in your mongocli profile. All values can also be
provided as flags or MCLI_* environment variables for scripted setup.`,
		Example: `  # Configure access to the DR state file interactively:
  mongocli ops-manager standby-clusters configure

  # Configure with static AWS credentials, no prompts:
  mongocli ops-manager standby-clusters configure \
    --s3BucketName my-dr-bucket \
    --s3Key my-cluster/dr_status_my-cluster.json \
    --awsRegion us-east-1 \
    --awsAuthMode staticCredentials \
    --awsAccessKey AKIA... --awsSecretKey ...`,
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			o.out = cmd.OutOrStdout()
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return o.Run(cmd.Context())
		},
	}

	addSharedFlags(cmd, &o.opts)

	return cmd
}
