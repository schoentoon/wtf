package githubnotifications

import "fmt"

func (widget *Widget) display() {
	widget.TextWidget.Redraw(widget.content)
}

func (widget *Widget) content() (string, string, bool) {
	// Choses the correct place to scroll to when changing sources
	if len(widget.View.GetHighlights()) > 0 {
		widget.View.ScrollToHighlight()
	} else {
		widget.View.ScrollToBeginning()
	}

	title := widget.CommonSettings().Title

	if widget.err != nil {
		return title, widget.err.Error(), false
	}

	var out string
	for idx, notification := range widget.Notifications {
		out += fmt.Sprintf(` ["%d"]%s %s in %s[""]`,
			idx,
			notification.GetSubject().GetType(),
			notification.GetSubject().GetTitle(),
			notification.GetRepository().GetFullName(),
		)
		out += "\n"
	}

	return title, out, false
}
