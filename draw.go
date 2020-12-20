package main

import (
	"fmt"
	"math"
	"syscall/js"
)

func drawLine(ctx js.Value, begin, end Vector, color string) {
	ctx.Set("strokeStyle", color)
	ctx.Call("beginPath")
	ctx.Call("moveTo", begin.X, begin.Y)
	ctx.Call("lineTo", end.X, end.Y)
	ctx.Call("closePath")
	ctx.Call("stroke")
}

func drawDot(ctx js.Value, pos Vector, color string) {
	ctx.Set("fillStyle", color)
	ctx.Call("beginPath")
	ctx.Call("arc", pos.X, pos.Y, 4, 0, 2*math.Pi, true)
	ctx.Call("fill")
}

func drawText(ctx js.Value, text string, x, y, sz int, color string) {
	ctx.Call("beginPath")
	ctx.Set("fillStyle", color)
	ctx.Set("font", fmt.Sprintf("bold %dpx Arial, sans-serif", sz))
	textWidth := ctx.Call("measureText", text).Get("width").Int()
	ctx.Call("fillText", text, x-(textWidth)/2, y)
}
