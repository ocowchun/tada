package command

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	git "gopkg.in/src-d/go-git.v4"

	"github.com/ocowchun/tada/utils"

	"github.com/mitchellh/cli"
)

type InstallCommand struct {
}

func (*InstallCommand) Help() string {
	return "Sorry, there're no help"
}

func (*InstallCommand) Run(args []string) int {
	packagePath := args[0]
	basePath, err := utils.ExpandHomeDir("~/.tada")
	if err != nil {
		log.Printf("%v", err)
	}
	// TODO validate packagePath
	url := "https://" + packagePath + ".git"
	strs := strings.Split(packagePath, "/")
	pluginName := strs[len(strs)-1]
	directory := basePath + "/plugins/" + pluginName
	fmt.Printf("git clone %s %s --recursive", url, directory)
	_, err = git.PlainClone(directory, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		log.Printf("Git clone with error: %v", err)
	}

	so := basePath + "/so/" + pluginName + ".so"
	mod := directory + "/main.go"
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
