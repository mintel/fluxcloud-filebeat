package config

type Config struct {
	FileBeatAddress string `envconfig:"filebeat_address" required:"false"`
	VCSRootURL      string `envconfig:"vcs_root_url" required:"true"`
	KeepFluxEvents  string `envconfig:"keep_flux_events" required:"false"`
	Port            int    `default:"8080"`
}
