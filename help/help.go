package help

import (
	"fmt"

	"github.com/olebedev/config"
	"github.com/wtfutil/wtf/app"
)

// Display displays the output of the --help argument
func Display(moduleName string, cfg *config.Config) {
	if moduleName == "" {
		fmt.Println("\n  --module takes a module name as an argument, i.e: '--module=github'")
	} else {
		fmt.Printf("%s\n", helpFor(moduleName, cfg))
	}
}

func helpFor(moduleName string, cfg *config.Config) string {
	cfg.Set("wtf.mods."+moduleName+".enabled", true)
	widget := app.MakeWidget(nil, nil, moduleName, cfg)

	// Since we are forcing enabled config, if no module
	// exists, we will get the unknown one
	if widget.CommonSettings().Title == "Unknown" {
		return "Unable to find module " + moduleName
	}

	result := ""
	result += utils.StripColorTags(widget.HelpText())
	result += "\n"
	result += "Configuration Attributes"
	result += widget.ConfigText()
	return result
}
