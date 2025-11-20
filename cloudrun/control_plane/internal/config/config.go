// Package config contains the configuration of the app.
package config

type Config struct {
	ProjectID            string   `env:"PROJECT_ID,required"`
	Region               string   `env:"REGION,required"`
	GithubAppPrivateKey  string   `env:"GITHUB_APP_PRIVATE_KEY,required"`
	InstanceTemplateName string   `env:"INSTANCE_TEMPLATE_NAME,required"`
	RunnerLabels         []string `env:"RUNNER_LABELS,required"`
	AllowedZones         []string `env:"ALLOWED_ZONES,required"`
	AppID                int64    `env:"GITHUB_APP_ID,required"`
	Port                 int      `env:"PORT,default=8080"`
	MaxRunnerCount       int      `env:"MAX_RUNNER_COUNT,required"`
	MinRunnerCount       int      `env:"MIN_RUNNER_COUNT,required"`
	UseStrictZoneOrder   bool     `env:"USE_STRICT_ZONE_ORDER,required"`
	Ephemeral            bool     `env:"EPHEMERAL,required"`
	UseJitConfig         bool     `env:"USE_JIT_CONFIG,required"`
	UseOrgRunners        bool     `env:"USE_ORG_RUNNERS,required"`
	Debug                bool     `env:"DEBUG,default=false"`
}
