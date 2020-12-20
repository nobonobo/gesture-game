package views

import (
	"gesture-game/components"
	"github.com/nobonobo/spago"
)

// Render ...
func (c *Playing) Render() spago.HTML {
	return spago.Tag("body", 
		spago.C(&components.Header{Title: "Playing", Abort: c.Abort()}),
		spago.Tag("canvas", 			
			spago.A("id", spago.S(`output`)),
		),
	)
}
