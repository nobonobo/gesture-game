package views

import (
	"gesture-game/components"
	"github.com/nobonobo/spago"
)

// Render ...
func (c *Top) Render() spago.HTML {
	return spago.Tag("body", 
		spago.C(&components.Header{Title: "ふうせんわりゲーム"}),
		spago.Tag("main", 			
			spago.A("class", spago.S(`container`)),
			spago.Tag("button", 				
				spago.Event("click", c.OnStart),
				spago.A("class", spago.S(`btn`)),
				spago.T(`すたーと！`),
			),
		),
	)
}
