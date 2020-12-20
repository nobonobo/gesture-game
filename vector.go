package main

import (
	"math"
	"syscall/js"
)

// Vector ...
type Vector struct {
	X float64
	Y float64
}

// Sub ...
func (v Vector) Sub(ov Vector) Vector { return Vector{v.X - ov.X, v.Y - ov.Y} }

// Norm ...
func (v Vector) Norm() float64 { return math.Sqrt(v.Dot(v)) }

// Dot ...
func (v Vector) Dot(ov Vector) float64 { return v.X*ov.X + v.Y*ov.Y }

// Distance ...
func (v Vector) Distance(ov Vector) float64 { return v.Sub(ov).Norm() }

// NewVector ...
func NewVector(v js.Value) Vector {
	return Vector{
		X: 640.0 - (v.Get("x").Float() * 640.0 / 300.0),
		Y: v.Get("y").Float() * 480.0 / 300.0,
	}
}
