package views

import (
	"gesture-game/components"
	"github.com/nobonobo/spago"
)

// Render ...
func (c *Result) Render() spago.HTML {
	return spago.Tag("body", 
		spago.C(&components.Header{Title: "けっか"}),
		spago.Tag("main", 			
			spago.A("class", spago.S(`container`)),
			spago.Tag("h2", 
				spago.T(``, spago.S(c.Score), `てん`),
			),
			spago.Tag("button", 				
				spago.Event("click", c.OnReturn),
				spago.A("class", spago.S(`btn`)),
				spago.T(`タイトルにもどる`),
			),
		),
	)
}
