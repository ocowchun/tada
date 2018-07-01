package command

import (
	"log"
	"os/exec"
	"strings"

	git "github.com/libgit2/git2go"
	"github.com/mitchellh/cli"
)

type InstallCommand struct {
}

func (*InstallCommand) Help() string {
	return "Sorry, there're no help"
}

func (*InstallCommand) Run(args []string) int {
	// fmt.Println(args)
	// url := "https://github.com/ocowchun/tada-hello-plugin.git"
	// "github.com/ocowchun/tada-hello-plugin"
	packagePath := args[0]
	strs := strings.Split(packagePath, "/")
	pluginName := strs[len(strs)-1]
	url := "https://" + packagePath + ".git"
	path := "/Users/ocowchun/go/src/github.com/ocowchun/tada/_plugins/tada-hello-plugin"
	options := &git.CloneOptions{Bare: false}
	_, err := git.Clone(url, path, options)
	if err != nil {
		log.Printf("Git clone with error: %v", err)
	}

	so := "/Users/ocowchun/go/src/github.com/ocowchun/tada/_plugins/.so/" + pluginName + ".so"
	mod := path + "/main.go"
	// buildCommand := "-buildmode=plugin -o " + so + " " + mod
	// buildCommand = "which go"
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", so, mod)
	log.Printf("Running command and waiting for it to finish...")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}
	log.Printf("Command finished with output: %v", string(output))
	return 0
}

func (*InstallCommand) Synopsis() string {
	return "Install tada package"
}

func InstallCommandFactory() (cli.Command, error) {
	return &InstallCommand{}, nil
}
