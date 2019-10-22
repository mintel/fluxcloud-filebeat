package config

type Config struct {
	FileBeatAddress string `envconfig:"filebeat_address" required:"false"`
	ClusterName     string `envconfig:"cluster_name" required:"true"`
	ClusterEnv      string `envconfig:"cluster_env" required:"true"`
	VCSRootURL      string `envconfig:"vcs_root_url" required:"true"`
	KeepFluxEvents  string `envconfig:"keep_flux_events" required:"false"`
	Port            int    `default:"8080"`
}
