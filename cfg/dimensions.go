package cfg

import (
	"github.com/olebedev/config"
	"github.com/wtfutil/wtf/utils"
)

// CalculateDimensions reads the module dimensions from the module and global config. The border is already substracted.
func CalculateDimensions(moduleConfig, globalConfig *config.Config) (int, int) {
	// Read the source data from the config
	left := moduleConfig.UInt("position.left", 0)
	top := moduleConfig.UInt("position.top", 0)
	width := moduleConfig.UInt("position.width", 0)
	height := moduleConfig.UInt("position.height", 0)

	cols := utils.ToInts(globalConfig.UList("wtf.grid.columns"))
	rows := utils.ToInts(globalConfig.UList("wtf.grid.rows"))

	// Make sure the values are in bounds
	left = utils.Clamp(left, 0, len(cols)-1)
	top = utils.Clamp(top, 0, len(rows)-1)
	width = utils.Clamp(width, 0, len(cols)-left)
	height = utils.Clamp(height, 0, len(rows)-top)

	// Start with the border subtracted and add all the spanned rows and cols
	w, h := -2, -2
	for _, x := range cols[left : left+width] {
		w += x
	}
	for _, y := range rows[top : top+height] {
		h += y
	}

	// The usable space may be empty
	w = utils.MaxInt(w, 0)
	h = utils.MaxInt(h, 0)

	return w, h
}
