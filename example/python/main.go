package main

import (
	"context"
	"flag"
	"fmt"
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
print(input())
print(input())
`
	input     = "1\ninput\n"
	filename  = "main.py"
	inputFile = "input.txt"

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

	files := []file.File{
		file.NewFile(filename, file.StringContent(code)),
		file.NewFile(inputFile, file.StringContent(input)),
	}

	client, err := coderunner.NewRunner(cfg, client.WithHost(dockerHost))
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	cmd := fmt.Sprintf(`cat %s | python3 %s`, inputFile, filename)
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
