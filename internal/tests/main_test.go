//go:build integration

package tests

import (
	"github.com/radium-rtf/coderunner_lib/config"
	r "github.com/radium-rtf/coderunner_lib/internal/tests/runner"
	"github.com/radium-rtf/coderunner_lib/limit"
	p "github.com/radium-rtf/coderunner_lib/profile"
	"os"
	"testing"
)

var (
	pythonProfile = p.Profile{
		Name:  "python3",
		Image: "python-coderunner:latest",
	}

	limits = limit.NewLimits()

	runner r.Runner
)

func TestMain(m *testing.M) {
	cfg := config.NewConfig(
		config.WithUser("sandbox"),
		config.WithUID(3456),
		config.WithWorkDir("/sandbox"),
	)

	runner = r.New(cfg)
	defer runner.Close()

	runner.SetUp()
	code := m.Run()
	runner.TearDown()

	os.Exit(code)
}
