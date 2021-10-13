package grafana_alert

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	widget "github.com/ocowchun/tada/Widget"
	"log"
	"os/exec"
	"sync"
)

type GrafanaAlert struct {
	previousKey string
	widgets.List
	alerts    []Alert
	alertLock sync.Mutex
	renderCh  chan<- struct{}
	apiKey    string
	hostName  string
}

func (l *GrafanaAlert) HandleUIEvent(e ui.Event) {
	switch e.ID {
	case "j", "<Down>":
		l.ScrollDown()
	case "k", "<Up>":
		l.ScrollUp()
	case "<C-d>":
		l.ScrollHalfPageDown()
	case "<C-u>":
		l.ScrollHalfPageUp()
	case "<C-f>":
		l.ScrollPageDown()
	case "<C-b>":
		l.ScrollPageUp()
	case "g":
		if l.previousKey == "g" {
			l.ScrollTop()
		}
	case "<Home>":
		l.ScrollTop()
	case "G", "<End>":
		l.ScrollBottom()
	case "<Enter>":
		story := l.selectedAlert()
		if story != nil {
			cmd := exec.Command("open", story.url)
			if _, err := cmd.Output(); err != nil {
				log.Println(err)
			}
		}
	}

	if l.previousKey == "g" {
		l.previousKey = ""
	} else {
		l.previousKey = e.ID
	}
}

func (l *GrafanaAlert) SetBorderStyle(style ui.Style) {
	l.BorderStyle = style
}

type Alert struct {
	title string
	url   string
}

func toRows(alerts []Alert) []string {
	result := make([]string, len(alerts))

	for i, alert := range alerts {
		result[i] = fmt.Sprintf(" %v", alert.title)
	}
	return result
}

func (l *GrafanaAlert) fetchAndUpdateAlerts() {
	alerts, err := fetchAlerts(l.hostName, l.apiKey)
	if err != nil {
		l.Rows = []string{err.Error()}
		return
	}

	l.updateAlerts(alerts)
}

func (l *GrafanaAlert) Refresh() {
	l.updateAlerts(make([]Alert, 0))
	l.Rows = []string{"loading..."}

	go l.fetchAndUpdateAlerts()
}

func (l *GrafanaAlert) updateAlerts(stories []Alert) {
	l.alertLock.Lock()
	defer l.alertLock.Unlock()

	l.alerts = stories
	l.Rows = toRows(stories)
}

func (l *GrafanaAlert) selectedAlert() *Alert {
	l.alertLock.Lock()
	defer l.alertLock.Unlock()

	i := l.SelectedRow
	if i >= len(l.alerts) {
		return nil
	}
	return &l.alerts[i]
}

func NewGrafanaAlert(config map[string]interface{}, renderCh chan<- struct{}) widget.Widget {
	l := widgets.NewList()
	l.Title = "Grafana Alert"
	l.SelectedRowStyle = ui.NewStyle(ui.ColorWhite, ui.ColorBlack, ui.ModifierUnderline)
	l.WrapText = false

	if _, existed := config["GRAFANA_HOSTNAME"]; !existed {
		log.Fatal("You must provide GRAFANA_HOSTNAME to enable this Widget")
	}
	if _, existed := config["GRAFANA_API_KEY"]; !existed {
		log.Fatal("You must provide GRAFANA_API_KEY to enable this Widget")
	}

	hostName := config["GRAFANA_HOSTNAME"].(string)
	apiKey := config["GRAFANA_API_KEY"].(string)

	component := &GrafanaAlert{
		List:     *l,
		renderCh: renderCh,
		hostName: hostName,
		apiKey:   apiKey,
	}
	go component.Refresh()
	return component
}
