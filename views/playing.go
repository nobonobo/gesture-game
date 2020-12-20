package views

import (
	"gesture-game/actions"
	"syscall/js"

	"github.com/nobonobo/spago"
	"github.com/nobonobo/spago/dispatcher"
)

//go:generate spago generate -c Playing -p views playing.html

// Playing  ...
type Playing struct {
	spago.Core
}

// Abort ...
func (c *Playing) Abort() spago.Markup {
	return spago.Tag("a",
		spago.ClassMap{"btn": true},
		spago.Event("click", c.OnAbort),
		spago.T("Abort"),
	)
}

// OnAbort ...
func (c *Playing) OnAbort(ev js.Value) {
	dispatcher.Dispatch(actions.GameEnd)
}
