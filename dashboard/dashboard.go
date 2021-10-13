package dashboard

import (
	"errors"
	ui "github.com/gizak/termui/v3"
	widget "github.com/ocowchun/tada/Widget"
	"github.com/ocowchun/tada/widgets"
	"github.com/ocowchun/tada/widgets/ghpr"
	"github.com/ocowchun/tada/widgets/grafana_alert"
	"github.com/ocowchun/tada/widgets/hackernews"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

type Dashboard struct {
}

var (
	defaultBorderStyle  = ui.NewStyle(ui.ColorWhite)
	selectedBorderStyle = ui.NewStyle(ui.ColorYellow)
	focusedBorderStyle  = ui.NewStyle(ui.ColorGreen)
)

type WidgetFactory func(map[string]interface{}, chan<- struct{}) widget.Widget

func Home() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	if currentUser.HomeDir == "" {
		return "", errors.New("cannot find user-specific home dir")
	}

	return currentUser.HomeDir, nil
}

func ExpandHomeDir(path string) (string, error) {
	if len(path) == 0 {
		return path, nil
	}

	if path[0] != '~' {
		return path, nil
	}

	if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
		return "", errors.New("cannot expand user-specific home dir")
	}

	dir, err := Home()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, path[1:]), nil
}

func findBasePath() string {
	basePath := os.Getenv("TADA_PATH")
	if basePath == "" {
		path, err := ExpandHomeDir("~/.tada")
		if err != nil {
			log.Printf("%v", err)
		}
		basePath = path
	}

	return basePath
}

type WidgetConfig struct {
	Name    string
	Width   int
	Height  int
	X       int
	Y       int
	Options map[string]interface{}
}

type Config struct {
	Widgets []WidgetConfig
}

func loadConfig() Config {
	basePath := findBasePath()
	raw, err := ioutil.ReadFile(basePath + "/tada.toml")
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = toml.Unmarshal(raw, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func (d Dashboard) Run() error {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	factoryMap := make(map[string]WidgetFactory)
	factoryMap["ghpr"] = ghpr.NewGitHubPRList
	factoryMap["myList"] = widgets.NewList
	factoryMap["hackernews"] = hackernews.NewHackerNews
	factoryMap["grafanaAlert"] = grafana_alert.NewGrafanaAlert

	ch := make(chan struct{})
	ready := make(chan struct{})
	ds := make([]ui.Drawable, 0)
	go func() {
		isReady := false
		for {
			select {
			case <-ch:
				if isReady {
					ui.Render(ds...)
				}
			case <-ready:
				isReady = true
			}
		}
	}()
	ch <- struct{}{}

	dashboardConfig := loadConfig()
	ws := make([]widget.Widget, 0)
	for _, widgetConfig := range dashboardConfig.Widgets {
		w := factoryMap[widgetConfig.Name](widgetConfig.Options, ch)
		w.SetRect(widgetConfig.X, widgetConfig.Y, widgetConfig.X+widgetConfig.Width, widgetConfig.Y+widgetConfig.Height)
		ws = append(ws, w)
	}
	for _, w := range ws {
		ds = append(ds, w)
	}
	ready <- struct{}{}

	var currentWidget widget.Widget
	currentWidgetIdx := -1
	widgetMode := false
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "<Escape>":
			return nil
		}

		if !widgetMode {
			if currentWidget != nil {
				currentWidget.SetBorderStyle(defaultBorderStyle)
			}
			switch e.ID {
			case "<Right>":
				currentWidgetIdx++
				if currentWidgetIdx == len(ws) {
					currentWidgetIdx = 0
				}
				currentWidget = ws[currentWidgetIdx]
				currentWidget.SetBorderStyle(selectedBorderStyle)
			case "<Left>":
				currentWidgetIdx--
				if currentWidgetIdx < 0 {
					currentWidgetIdx = len(ws) - 1
				}
				currentWidget = ws[currentWidgetIdx]
				currentWidget.SetBorderStyle(selectedBorderStyle)
			case "<Enter>":
				if currentWidget != nil {
					widgetMode = true
					currentWidget.SetBorderStyle(focusedBorderStyle)
				}
			case "r":
				for _, w := range ws {
					w.Refresh()
				}
			}

		} else {
			if e.ID == "q" {
				currentWidget.SetBorderStyle(selectedBorderStyle)
				widgetMode = false
			} else {
				currentWidget.HandleUIEvent(e)
			}
		}
		ch <- struct{}{}

	}
	return nil
}

func New() *Dashboard {
	return &Dashboard{}
}
