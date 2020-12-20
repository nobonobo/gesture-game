package views

import (
	"github.com/nobonobo/spago"
)

//go:generate spago generate -c Playing -p views playing.html

// Playing  ...
type Playing struct {
	spago.Core
}
