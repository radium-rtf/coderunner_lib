package config

type Opt func(*Config)

func WithUser(user string) Opt {
	return func(config *Config) {
		config.User = user
	}
}

func WithUID(uid int) Opt {
	return func(config *Config) {
		config.Uid = uid
	}
}

func WithWorkDir(workDir string) Opt {
	return func(config *Config) {
		config.Workdir = workDir
	}
}
