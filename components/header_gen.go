package components

import (
	"github.com/nobonobo/spago"
)

// Render ...
func (c *Header) Render() spago.HTML {
	return spago.Tag("header", 		
		spago.A("class", spago.S(`navbar`)),
		spago.Tag("section", 			
			spago.A("class", spago.S(`navber-section`)),
			spago.Tag("a", 				
				spago.A("class", spago.S(`navbar-brand mr-2`)),
				spago.T(``, spago.S(c.Title), ``),
			),
		),
		spago.Tag("section", 			
			spago.A("class", spago.S(`navber-section`)),
			spago.If(c.Abort!=nil, c.Abort),
		),
	)
}
