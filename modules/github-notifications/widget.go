package githubnotifications

import (
	"context"
	"net/http"
	"strconv"

	"github.com/google/go-github/v26/github"
	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/logger"
	"github.com/wtfutil/wtf/utils"
	"github.com/wtfutil/wtf/view"
	"golang.org/x/oauth2"
)

// Widget define wtf widget to register widget later
type Widget struct {
	view.KeyboardWidget
	view.ScrollableWidget

	oauthClient   *http.Client
	Client        *github.Client
	Notifications []*github.Notification
	err           error

	settings *Settings
	Selected int
}

// NewWidget creates a new instance of the widget
func NewWidget(app *tview.Application, pages *tview.Pages, settings *Settings) *Widget {
	widget := Widget{
		KeyboardWidget:   view.NewKeyboardWidget(app, pages, settings.common),
		ScrollableWidget: view.NewScrollableWidget(app, settings.common),

		settings: settings,
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: settings.apiKey})
	widget.oauthClient = oauth2.NewClient(context.Background(), ts)

	widget.Client = github.NewClient(widget.oauthClient)

	widget.initializeKeyboardControls()
	widget.View.SetRegions(true)
	widget.View.SetInputCapture(widget.InputCapture)

	widget.Unselect()

	widget.KeyboardWidget.SetView(widget.View)

	return &widget
}

/* -------------------- Exported Functions -------------------- */

// GetSelected returns the index of the currently highlighted item as an int
func (widget *Widget) GetSelected() int {
	if widget.Selected < 0 {
		return 0
	}
	return widget.Selected
}

// Next cycles the currently highlighted text down
func (widget *Widget) Next() {
	widget.Selected++
	if widget.Selected >= len(widget.Notifications) {
		widget.Selected = 0
	}
	widget.View.Highlight(strconv.Itoa(widget.Selected)).ScrollToHighlight()
}

// Prev cycles the currently highlighted text up
func (widget *Widget) Prev() {
	widget.Selected--
	if widget.Selected < 0 {
		widget.Selected = len(widget.Notifications) - 1
	}
	widget.View.Highlight(strconv.Itoa(widget.Selected)).ScrollToHighlight()
}

// Unselect stops highlighting the text and jumps the scroll position to the top
func (widget *Widget) Unselect() {
	widget.Selected = -1
	widget.View.Highlight()
	widget.View.ScrollToBeginning()
}

// Refresh reloads the github data via the Github API and reruns the display
func (widget *Widget) Refresh() {
	notifications, _, err := widget.Client.Activity.ListNotifications(context.Background(), &github.NotificationListOptions{All: true})
	if err != nil {
		widget.err = err
		widget.Notifications = nil
	} else {
		widget.err = nil
		widget.Notifications = notifications
	}
	widget.display()
}

// HelpText displays the widgets controls
func (widget *Widget) HelpText() string {
	return widget.KeyboardWidget.HelpText()
}

/* -------------------- Unexported Functions -------------------- */

type HtmlResponse struct {
	HtmlURL string `json:"html_url"`
}

func (widget *Widget) openNotification() {
	currentSelection := widget.View.GetHighlights()
	if widget.Selected >= 0 && currentSelection[0] != "" {
		url := widget.Notifications[widget.GetSelected()].GetSubject().GetURL()
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			logger.Log(err.Error())
			return
		}

		res, err := widget.oauthClient.Do(req)
		if err != nil {
			logger.Log(err.Error())
			return
		}
		defer res.Body.Close()

		var resp HtmlResponse
		err = utils.ParseJSON(&resp, res.Body)
		if err != nil {
			logger.Log(err.Error())
			return
		}
		utils.OpenFile(resp.HtmlURL)
	}
}
