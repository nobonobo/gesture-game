package components

import (
	"github.com/nobonobo/spago"
)

//go:generate spago generate -c Header -p components header.html

// Header  ...
type Header struct {
	spago.Core
	Title string
	Abort spago.Markup
}
