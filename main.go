package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

const (
	VERSION = "0.0.1"
)

// git clone ssh://git@git.maujagroup.com:1960/ms/mercadolibre-webhook.git `go env GOPATH`/src/git.maujagroup.com/ms/mercadolibre-webhook
func main() {
	flag.Parse()
	repoStr := strings.Trim(flag.Arg(0), "")

	if repoStr == "version" {
		fmt.Println(VERSION)
		return
	}
	repo, err := url.Parse(repoStr)
	if err != nil {
		log.Fatal(err)
	}
	path := strings.Replace(".git", "", repo.Path, -1)

	gopath, err := getOutput("go env GOPATH")
	if err != nil {
		log.Fatal(err)
		return
	}
	project := fmt.Sprintf("%s/src/%s%s",
		*gopath, strings.Trim(repo.Hostname(), " "), path,
	)
	log.Println("Clonando en:", project)

	cmdClone := fmt.Sprintf("git clone %s %s", repoStr, project)
	// log.Println("RUN en:", cmdClone)
	err = run(cmdClone)
	if err != nil {
		log.Println("::err::")
		log.Fatal(err)
		return
	}
}

func getOutput(cmdStr string) (*string, error) {
	cmdSplit := strings.Split(cmdStr, " ")
	if len(cmdSplit) == 0 {
		return nil, fmt.Errorf("cmd is required")
	}
	args := []string{}
	for i := 1; i < len(cmdSplit); i++ {
		args = append(args, cmdSplit[i])
	}
	cmd := exec.Command(cmdSplit[0], args...)
	stdout, err := cmd.Output()

	if err != nil {
		return nil, err
	}
	out := strings.Trim(string(stdout), "\n")
	return &out, nil
}

func run(cmdStr string) error {
	cmdSplit := strings.Split(cmdStr, " ")
	if len(cmdSplit) == 0 {
		return fmt.Errorf("cmd is required")
	}
	args := []string{}
	for i := 1; i < len(cmdSplit); i++ {
		args = append(args, cmdSplit[i])
	}
	cmd := exec.Command(cmdSplit[0], args...)
	cmd.Stdout = os.Stdout // or any other io.Writer
	cmd.Stderr = os.Stderr // or any other io.Writer
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
