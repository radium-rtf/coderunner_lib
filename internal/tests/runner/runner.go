package runner

import (
	"fmt"
	"github.com/docker/docker/client"
	coderunner "github.com/radium-rtf/coderunner_lib"
	"github.com/radium-rtf/coderunner_lib/config"
	"os"
)

const (
	dockerHostKey = "TEST_DOCKER_HOST"
)

type Runner struct {
	dockerHost string
	*coderunner.Runner
}

func New(cfg *config.Config) Runner {
	dockerHost := os.Getenv(dockerHostKey)
	if dockerHost == "" {
		panic(fmt.Sprintf("key isnt  %s set", dockerHostKey))
	}

	runner, err := coderunner.NewRunner(cfg, client.WithHost(dockerHost))
	if err != nil {
		panic(err)
	}
	return Runner{dockerHost: dockerHost, Runner: runner}
}
