package command

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/cli"
	"github.com/ocowchun/tada/dashboard"
	"github.com/ocowchun/tada/utils"
	"github.com/ocowchun/tada/widget"
)

type InitCommand struct {
}

func (*InitCommand) Help() string {
	return "Initialize tada required config"
}

func (*InitCommand) Run(args []string) int {
	basePath := utils.FindBasePath()
	if err := os.MkdirAll(basePath, 0755); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	path := basePath + "/tada.toml"
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	defer f.Close()

	opts := make(map[string]interface{})
	opts["GITHUB_USERNAME"] = "your-github-username"
	opts["GITHUB_TOKEN"] = "your-github-developer-token"
	w := widget.Config{
		Name:    "tada-github-pr",
		Title:   "Pull Requests",
		Width:   3,
		Height:  3,
		X:       1,
		Y:       0,
		Options: opts,
	}
	config := dashboard.Config{
		Widgets: []widget.Config{w},
	}

	writer := bufio.NewWriter(f)
	if err = toml.NewEncoder(writer).Encode(config); err != nil {
		log.Fatal(err)
	}
	writer.Flush()

	fmt.Println("Create tada.toml in " + basePath + " successfully")
	return 0
}

func (*InitCommand) Synopsis() string {
	return "Initialize tada required config"
}

func InitCommandFactory() (cli.Command, error) {
	return &InitCommand{}, nil
}
