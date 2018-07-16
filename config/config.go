package dashboard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type WidgetConfig struct {
	Name    string                 `json:"name"`
	Width   int                    `json:"width"`
	Height  int                    `json:"height"`
	X       int                    `json:"x"`
	Y       int                    `json:"y"`
	Options map[string]interface{} `json:"options"`
}

type Config struct {
	Widgets []WidgetConfig `json:""widgets""`
}

func LoadConfig(path string) Config {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var c Config
	json.Unmarshal(raw, &c)
	return c
}
