package views

import (
	"github.com/nobonobo/spago"
)

// Render ...
func (c *Playing) Render() spago.HTML {
	return spago.Tag("body", 
		spago.Tag("canvas", 			
			spago.A("id", spago.S(`output`)),
		),
		
		spago.Tag("img", 			
			spago.A("id", spago.S(`baloon`)),
			spago.A("src", spago.S(`assets/baloon.png`)),
		),
		
		spago.Tag("img", 			
			spago.A("id", spago.S(`spark`)),
			spago.A("src", spago.S(`assets/spark.png`)),
		),
		spago.Tag("audio", 			
			spago.A("id", spago.S(`clap`)),
			spago.A("src", spago.S(`assets/clap.mp3`)),
		),
	)
}
