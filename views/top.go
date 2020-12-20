package views

import (
	"gesture-game/actions"
	"syscall/js"

	"github.com/nobonobo/spago"
	"github.com/nobonobo/spago/dispatcher"
)

//go:generate spago generate -c Top -p views top.html

// Top  ...
type Top struct {
	spago.Core
}

// OnStart ...
func (c *Top) OnStart(ev js.Value) {
	dispatcher.Dispatch(actions.GameStart)
}
