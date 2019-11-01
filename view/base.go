package view

import (
	"fmt"
	"sync"

	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/cfg"
	"github.com/wtfutil/wtf/utils"
	"github.com/wtfutil/wtf/help"
)

type Base struct {
	app             *tview.Application
	bordered        bool
	commonSettings  *cfg.Common
	enabled         bool
	focusChar       string
	focusable       bool
	name            string
	quitChan        chan bool
	refreshing      bool
	refreshInterval int
	enabledMutex    *sync.Mutex
}

func NewBase(app *tview.Application, commonSettings *cfg.Common) Base {
	base := Base{
		commonSettings:  commonSettings,
		app:             app,
		bordered:        commonSettings.Bordered,
		enabled:         commonSettings.Enabled,
		focusChar:       commonSettings.FocusChar(),
		focusable:       commonSettings.Focusable,
		name:            commonSettings.Name,
		quitChan:        make(chan bool),
		refreshInterval: commonSettings.RefreshInterval,
		refreshing:      false,
		enabledMutex:    &sync.Mutex{},
	}
	return base
}

/* -------------------- Exported Functions -------------------- */

// Bordered returns whether or not this widget should be drawn with a border
func (base *Base) Bordered() bool {
	return base.bordered
}

func (base *Base) BorderColor() string {
	if base.Focusable() {
		return base.commonSettings.Colors.BorderFocusable
	}

	return base.commonSettings.Colors.BorderNormal
}

func (base *Base) CommonSettings() *cfg.Common {
	return base.commonSettings
}

func (base *Base) ConfigText() string {
	return help.HelpFromInterface(cfg.Common{})
}

func (base *Base) ContextualTitle(defaultStr string) string {
	if defaultStr == "" && base.FocusChar() == "" {
		return ""
	} else if defaultStr != "" && base.FocusChar() == "" {
		return fmt.Sprintf(" %s ", defaultStr)
	} else if defaultStr == "" && base.FocusChar() != "" {
		return fmt.Sprintf(" [darkgray::u]%s[::-][green] ", base.FocusChar())
	} else {
		return fmt.Sprintf(" %s [darkgray::u]%s[::-][green] ", defaultStr, base.FocusChar())
	}
}

func (base *Base) Disable() {
	base.enabledMutex.Lock()
	base.enabled = false
	base.enabledMutex.Unlock()
}

func (base *Base) Disabled() bool {
	base.enabledMutex.Lock()
	result := !base.enabled
	base.enabledMutex.Unlock()
	return result
}

func (base *Base) Enabled() bool {
	base.enabledMutex.Lock()
	result := base.enabled
	base.enabledMutex.Unlock()
	return result
}

func (base *Base) Focusable() bool {
	base.enabledMutex.Lock()
	result := base.enabled && base.focusable
	base.enabledMutex.Unlock()
	return result
}

func (base *Base) FocusChar() string {
	return base.focusChar
}

func (base *Base) HelpText() string {
	return fmt.Sprintf("\n  There is no help available for widget %s", base.commonSettings.Module.Type)
}

func (base *Base) Name() string {
	return base.name
}

func (base *Base) QuitChan() chan bool {
	return base.quitChan
}

// Refreshing returns TRUE if the base is currently refreshing its data, FALSE if it is not
func (base *Base) Refreshing() bool {
	return base.refreshing
}

// RefreshInterval returns how often, in seconds, the base will return its data
func (base *Base) RefreshInterval() int {
	return base.refreshInterval
}

func (base *Base) SetFocusChar(char string) {
	base.focusChar = char
}

func (base *Base) Stop() {
	base.enabledMutex.Lock()
	base.enabled = false
	base.enabledMutex.Unlock()
	base.quitChan <- true
}

func (base *Base) String() string {
	return base.name
}
