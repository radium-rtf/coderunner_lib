package config

const (
	DockerWorkDir = "/sandbox"
	DefaultUser   = "sandbox"
)

type Config struct {
	User       string
	Uid        int
	Workdir    string
	DockerHost string
}

func NewConfig(opts ...Opt) *Config {
	cfg := &Config{
		Workdir:    "/sandbox",
		User:       "root",
		Uid:        3456,
		DockerHost: "unix:///var/run/docker.sock",
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}
