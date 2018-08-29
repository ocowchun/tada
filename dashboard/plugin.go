package dashboard

import (
	"fmt"
	"os"
	"plugin"

	"github.com/ocowchun/tada/utils"
	widget "github.com/ocowchun/tada/widget"
)

func LoadPlugin(pluginName string, config widget.Config) widget.Box {
	basePath := utils.FindBasePath()
	mod := basePath + "/so/" + pluginName + ".so"
	plug, err := plugin.Open(mod)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	symNewBox, err := plug.Lookup("NewBox")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	newBox, ok := symNewBox.(func(config widget.Config) widget.Box)
	if !ok {
		fmt.Println("unexpected type from module symbol")
		os.Exit(1)
	}

	return newBox(config)
}
