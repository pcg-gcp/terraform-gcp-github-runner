package config

type Config struct {
	ProjectID            string   `env:"PROJECT_ID,required"`
	Zone                 string   `env:"ZONE,required"`
	GithubAppPrivateKey  string   `env:"GITHUB_APP_PRIVATE_KEY,required"`
	InstanceTemplateName string   `env:"INSTANCE_TEMPLATE_NAME,required"`
	RunnerLabels         []string `env:"RUNNER_LABELS,required"`
	AppID                int64    `env:"GITHUB_APP_ID,required"`
	Port                 int      `env:"PORT,default=8080"`
	MaxRunnerCount       int      `env:"MAX_RUNNER_COUNT,required"`
	MinRunnerCount       int      `env:"MIN_RUNNER_COUNT,required"`
	Ephemeral            bool     `env:"EPHEMERAL,required"`
	UseJitConfig         bool     `env:"USE_JIT_CONFIG,required"`
	UseOrgRunners        bool     `env:"USE_ORG_RUNNERS,required"`
	Debug                bool     `env:"DEBUG,default=false"`
}
