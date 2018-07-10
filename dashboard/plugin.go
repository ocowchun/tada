package dashboard

import (
	"fmt"
	"os"
	"plugin"

	"github.com/ocowchun/tada/utils"
	widget "github.com/ocowchun/tada/widget"
)

func LoadPlugin(pluginName string) *widget.Widget {
	basePath := utils.FindBasePath()
	mod := basePath + "/so/" + pluginName + ".so"
	plug, err := plugin.Open(mod)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	symNewWidget, err := plug.Lookup("NewWidget")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	newWidget, ok := symNewWidget.(func() *widget.Widget)
	if !ok {
		fmt.Println("unexpected type from module symbol")
		os.Exit(1)
	}
	return newWidget()
}
