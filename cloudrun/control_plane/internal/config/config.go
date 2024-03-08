package config

type Config struct {
	ProjectID            string `env:"PROJECT_ID,required"`
	Zone                 string `env:"ZONE,required"`
	GithubAppPrivateKey  string `env:"GITHUB_APP_PRIVATE_KEY,required"`
	InstanceTemplateName string `env:"INSTANCE_TEMPLATE_NAME,required"`
	AppID                int64  `env:"GITHUB_APP_ID,required"`
	Port                 int    `env:"PORT,default=8080"`
	Debug                bool   `env:"DEBUG,default=false"`
}
