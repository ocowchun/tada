package main

import (
	"fmt"
	"log"
	"os"
	"plugin"

	"github.com/mitchellh/cli"
	command "github.com/ocowchun/tada/command"
)

func readInstalledPackages() []string {
	return []string{"hello"}
}

type TestCommand struct {
}

func (*TestCommand) Help() string {
	return "Sorry, there're no help"
}

type Greeter interface {
	Greet()
}

func (*TestCommand) Run(args []string) int {
	// packages := readInstalledPackages()
	fmt.Println("yoyo")
	mod := "./_plugins/.so/hello.so"
	plug, err := plugin.Open(mod)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	symGreeter, err := plug.Lookup("Greeter")
	if err != nil {
		fmt.Println(err)
		return 1
	}

	var greeter Greeter
	greeter, ok := symGreeter.(Greeter)
	if !ok {
		fmt.Println("unexpected type from module symbol")
		return 1
	}

	// 4. use the module
	greeter.Greet()

	return 0
}

func (*TestCommand) Synopsis() string {
	return "test"
}

func testCommandFactory() (cli.Command, error) {
	return &TestCommand{}, nil
}

func main() {
	c := cli.NewCLI("app", "0.0.1")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"run":     command.RunCommandFactory,
		"install": command.InstallCommandFactory,
		"link":    command.LinkCommandFactory,
		"test":    testCommandFactory,
		"init":    command.InitCommandFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)

}
