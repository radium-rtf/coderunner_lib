package main

import (
	"context"
	"flag"
	"github.com/docker/docker/client"
	coderunner "github.com/radium-rtf/coderunner_lib"
	"github.com/radium-rtf/coderunner_lib/config"
	"github.com/radium-rtf/coderunner_lib/file"
	"github.com/radium-rtf/coderunner_lib/limit"
	"github.com/radium-rtf/coderunner_lib/profile"
	"log"
	"os"
)

const (
	code = `
print(11)
print(111)
print(1111)
print(11111)
print(111111)
`
	filename = "main.py"

	timeoutInSec = 1
)

var (
	dockerHost string
)

func init() {
	set := flag.NewFlagSet("cli", flag.ContinueOnError)
	set.StringVar(&dockerHost, "host", "unix:///var/run/docker.sock", "--host=<str>")
	err := set.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}
}

func main() {
	ctx := context.Background()

	profile := profile.Profile{
		Name:  "python",
		Image: "python-coderunner:latest",
	}

	cfg := config.NewConfig(
		config.WithUser("sandbox"),
		config.WithUID(3456),
		config.WithWorkDir("/sandbox"),
	)

	limits := limit.NewLimits(limit.WithTimeoutInSec(timeoutInSec))

	files := []file.File{file.NewFile(filename, file.StringContent(code))}

	client, err := coderunner.NewRunner(cfg, client.WithHost(dockerHost))
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	cmd := []string{"python3", filename}
	sandbox, err := client.NewSandbox(ctx, cmd, profile, limits, files)
	if err != nil {
		log.Fatalln(err)
	}
	defer sandbox.Close()

	if err := sandbox.Start(); err != nil {
		log.Fatalln(err)
	}

	statusCode, err := sandbox.Wait()
	if err != nil {
		log.Println("err", err.Error())
	}
	log.Println("statusCode", statusCode)

	out, err := sandbox.ShowStd()
	if err != nil {
		log.Println("err", err.Error())
	}

	log.Println(out.StdOut)
	log.Println("stdErr:")
	log.Println(out.StdErr)
}
