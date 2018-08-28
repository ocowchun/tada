package dashboard

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

type Widget struct {
	Name    string
	Width   int
	Height  int
	X       int
	Y       int
	Options map[string]interface{}
}

type Config struct {
	Widgets []Widget
}

func LoadConfig(path string) Config {
	var config Config
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if _, err = toml.Decode(string(raw), &config); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return config
}
