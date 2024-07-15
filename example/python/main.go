package main

import (
	"context"
	"flag"
	"github.com/docker/docker/client"
	"github.com/khostya/coderunner"
	"github.com/khostya/coderunner/config"
	"github.com/khostya/coderunner/file"
	"github.com/khostya/coderunner/limit"
	"github.com/khostya/coderunner/profile"
	"log"
	"os"
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

	limits := limit.NewLimits()

	code := `
print(11)
print(111)
print(1111)
print(11111)
print(111111)
`
	files := []file.File{
		{Name: "main.py", Content: code},
	}

	client, err := coderunner.NewRunner(cfg, client.WithHost(dockerHost))
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	cmd := []string{"python3", "main.py"}
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

	out, err := sandbox.ShowStdout()
	if err != nil {
		log.Println("err", err.Error())
	}

	log.Println(out)
}
