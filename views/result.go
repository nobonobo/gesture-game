package views

import (
	"syscall/js"

	"github.com/nobonobo/spago"
	"github.com/nobonobo/spago/router"
)

//go:generate spago generate -c Result -p views result.html

// Result  ...
type Result struct {
	spago.Core
	Score int
}

// OnReturn ...
func (c *Result) OnReturn(ev js.Value) {
	router.Navigate("/")
}
