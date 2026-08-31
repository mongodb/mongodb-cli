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

package standby

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go/logging"
)

// AuthMode represents the authentication mode for AWS credentials.
type AuthMode string

const (
	// AuthModeAssumeRole uses an IAM role ARN for auto-refreshing credentials.
	AuthModeAssumeRole AuthMode = "assumeRole"
	// AuthModeStaticCredentials uses long-lived IAM user credentials (access key and secret key, no session token).
	AuthModeStaticCredentials AuthMode = "staticCredentials"
	// AuthModeCredentialsChain uses the default AWS credentials chain (env vars, ~/.aws/credentials, IAM role, etc.).
	AuthModeCredentialsChain AuthMode = "credentialsChain"
)

var validAuthModes = []AuthMode{AuthModeAssumeRole, AuthModeStaticCredentials, AuthModeCredentialsChain}

// ValidAuthModes returns the supported AWS authentication modes.
func ValidAuthModes() []AuthMode {
	return validAuthModes
}

// Credentials holds fields needed to create a s3 client.
type Credentials struct {
	AuthMode        AuthMode
	RoleArn         string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	BucketName      string
	Region          string
	AWSProfile      string
	Endpoint        string
	UsePathStyle    bool
}

// Validate validates the aws credentials based on the AuthMode.
//   - assumeRole: Requires RoleArn, BucketName, Region
//   - staticCredentials: Requires AccessKeyID, SecretAccessKey, BucketName, Region (no SessionToken)
//   - credentialsChain: Requires only BucketName, Region (AWS SDK default chain / instance profile)
func (c Credentials) Validate() error {
	if c.BucketName == "" {
		return errors.New("aws.bucketName must be provided")
	}
	if c.Region == "" {
		return errors.New("aws.region must be provided")
	}

	switch c.AuthMode {
	case AuthModeAssumeRole:
		if c.RoleArn == "" {
			return errors.New("aws.roleArn must be provided when aws.authMode is assumeRole")
		}
		if c.AccessKeyID != "" || c.SecretAccessKey != "" || c.SessionToken != "" {
			return errors.New("aws.accessKeyId, aws.secretAccessKey, and aws.sessionToken must not be provided when aws.authMode is assumeRole")
		}
	case AuthModeStaticCredentials:
		if c.AccessKeyID == "" {
			return errors.New("aws.accessKeyId must be provided when aws.authMode is staticCredentials")
		}
		if c.SecretAccessKey == "" {
			return errors.New("aws.secretAccessKey must be provided when aws.authMode is staticCredentials")
		}
		if c.RoleArn != "" {
			return errors.New("aws.roleArn must not be provided when aws.authMode is staticCredentials")
		}
		if c.SessionToken != "" {
			return errors.New("aws.sessionToken must not be provided when aws.authMode is staticCredentials")
		}
	case AuthModeCredentialsChain:
		if c.RoleArn != "" || c.AccessKeyID != "" || c.SecretAccessKey != "" || c.SessionToken != "" {
			return errors.New("aws.roleArn, aws.accessKeyId, aws.secretAccessKey, and aws.sessionToken must not be provided when aws.authMode is credentialsChain")
		}
	default:
		return fmt.Errorf("aws.authMode must be one of %v, got %s", validAuthModes, c.AuthMode)
	}

	return nil
}

func newRawS3Client(ctx context.Context, creds Credentials, extraOpts ...func(*s3.Options)) (*s3.Client, error) {
	if err := creds.Validate(); err != nil {
		return nil, err
	}

	switch creds.AuthMode {
	case AuthModeAssumeRole:
		return createClientWithAssumeRoleProvider(ctx, creds, extraOpts...)
	case AuthModeStaticCredentials:
		return createClientWithCredsProvider(creds, extraOpts...)
	case AuthModeCredentialsChain:
		return createClientWithDefaultChain(ctx, creds, extraOpts...)
	default:
		return nil, fmt.Errorf("invalid aws.authMode: %s, must be one of %v", creds.AuthMode, validAuthModes)
	}
}

func loadChainOptions(creds Credentials) []func(*awsconfig.LoadOptions) error {
	opts := []func(*awsconfig.LoadOptions) error{awsconfig.WithRegion(creds.Region)}
	if creds.AWSProfile != "" {
		opts = append(opts, awsconfig.WithSharedConfigProfile(creds.AWSProfile))
	}
	return opts
}

func createClientWithDefaultChain(
	ctx context.Context,
	creds Credentials,
	extraOpts ...func(*s3.Options),
) (*s3.Client, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx, loadChainOptions(creds)...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config for credentialsChain: %w", err)
	}
	return s3.NewFromConfig(cfg, mergeS3Options(creds, extraOpts)...), nil
}

func createClientWithAssumeRoleProvider(
	ctx context.Context,
	creds Credentials,
	extraOpts ...func(*s3.Options),
) (*s3.Client, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx, loadChainOptions(creds)...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config for assumeRole: %w", err)
	}
	stsClient := sts.NewFromConfig(cfg)
	assumeRoleProvider := stscreds.NewAssumeRoleProvider(stsClient, creds.RoleArn, func(options *stscreds.AssumeRoleOptions) {
		options.RoleSessionName = "mongocli-standby-clusters"
	})

	cfg.Credentials = aws.NewCredentialsCache(assumeRoleProvider)
	return s3.NewFromConfig(cfg, mergeS3Options(creds, extraOpts)...), nil
}

func createClientWithCredsProvider(
	creds Credentials,
	extraOpts ...func(*s3.Options),
) (*s3.Client, error) {
	cfg := aws.Config{
		Region: creds.Region,
		Credentials: credentials.NewStaticCredentialsProvider(
			creds.AccessKeyID,
			creds.SecretAccessKey,
			creds.SessionToken,
		),
	}
	return s3.NewFromConfig(cfg, mergeS3Options(creds, extraOpts)...), nil
}

func s3Options(creds Credentials) func(*s3.Options) {
	return func(o *s3.Options) {
		if creds.Endpoint != "" {
			o.BaseEndpoint = &creds.Endpoint
		}
		o.UsePathStyle = creds.UsePathStyle
		o.Logger = logging.Nop{}
	}
}

func mergeS3Options(creds Credentials, extraOpts []func(*s3.Options)) []func(*s3.Options) {
	opts := []func(*s3.Options){s3Options(creds)}
	return append(opts, extraOpts...)
}
