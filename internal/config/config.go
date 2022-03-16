// This package contains the general configuration for repositories
// and for the lambda in general
package config

import (
	"context"

	"github.com/cyralinc/sidecar-failopen/internal/secrets"
	"github.com/spf13/viper"
)

type PostgreSQLConfig struct {
	ConnectionStringOptions string
}

type SnowflakeConfig struct {
	Account   string
	Role      string
	Warehouse string
}

// RepoConfig is the configuration for a repository, including
// connection information and metadata.
type RepoConfig struct {
	Host                    string
	Port                    int
	User                    string
	Password                string
	Database                string
	RepoName                string
	RepoType                string
	ConnectionStringOptions string
	SnowflakeConfig         SnowflakeConfig
	ConnectionTimeout       int
}

// LambdaConfig is the configuration for the lambda in general. This
// struct is to be used via the global Config function.
type LambdaConfig struct {
	NumberOfRetries int
	LogLevel        string
	StackName       string
	Sidecar         RepoConfig
	Repo            RepoConfig
}

func init() {
	// using prefix to get all environment variables starting with
	// "FAIL_OPEN" as configuration entries on viper
	viper.SetEnvPrefix("FAIL_OPEN")

	// sidecar location configuration
	viper.BindEnv("sidecar_port")
	viper.BindEnv("sidecar_host")
	viper.BindEnv("sidecar_timeout")

	// repository configuration
	viper.BindEnv("repo_type")
	viper.BindEnv("repo_host")
	viper.BindEnv("repo_port")
	viper.BindEnv("repo_name")
	viper.BindEnv("repo_database")
	viper.BindEnv("repo_secret")
	viper.BindEnv("repo_timeout")

	viper.BindEnv("n_retries") // number of retries on each healthcheck
	viper.BindEnv("log_level") // log level for the lambda

	viper.BindEnv("cf_stack_name") // name of the stack

	viper.BindEnv("connection_string_options") // connection options for pg based repos

	// snowflake configuration
	viper.BindEnv("snowflake_role")
	viper.BindEnv("snowflake_account")
	viper.BindEnv("snowflake_warehouse")
}

var c *LambdaConfig

// Config returns the global configuration for the lambda function, initializing
// all values and recovering the secrets.
func Config() *LambdaConfig {
	secret := viper.GetString("repo_secret")
	if c == nil {
		sec, err := secrets.RepoSecretFromSecretsManager(context.Background(), secret)
		if err != nil {
			panic(err)
		}

		c = &LambdaConfig{
			NumberOfRetries: viper.GetInt("n_retries"),
			LogLevel:        viper.GetString("log_level"),
			Repo: RepoConfig{
				Host:                    viper.GetString("repo_host"),
				Port:                    viper.GetInt("repo_port"),
				Database:                viper.GetString("repo_database"),
				RepoType:                viper.GetString("repo_type"),
				RepoName:                viper.GetString("repo_name"),
				User:                    sec.Username,
				Password:                sec.Password,
				ConnectionStringOptions: viper.GetString("pg_conn_opts"),

				SnowflakeConfig: SnowflakeConfig{
					Account:   viper.GetString("snowflake_account"),
					Role:      viper.GetString("snowflake_role"),
					Warehouse: viper.GetString("snowflake_warehouse"),
				},
				ConnectionTimeout: viper.GetInt("repo_timeout"),
			},
			Sidecar: RepoConfig{
				Host:                    viper.GetString("sidecar_host"),
				Port:                    viper.GetInt("sidecar_port"),
				Database:                viper.GetString("repo_database"),
				RepoType:                viper.GetString("repo_type"),
				RepoName:                viper.GetString("repo_name"),
				User:                    sec.Username,
				Password:                sec.Password,
				ConnectionStringOptions: viper.GetString("connection_string_options"),
				SnowflakeConfig: SnowflakeConfig{
					Account:   viper.GetString("snowflake_account"),
					Role:      viper.GetString("snowflake_role"),
					Warehouse: viper.GetString("snowflake_warehouse"),
				},
				ConnectionTimeout: viper.GetInt("sidecar_timeout"),
			},
			StackName: viper.GetString("cf_stack_name"),
		}
	}

	return c
}
