package main

import (
	"log"
	"syscall/js"
	"time"
)

var (
	baloon, spark js.Value
	// Baloons ...
	Baloons = map[*Baloon]bool{}
)

// Baloon ...
type Baloon struct {
	pos      Vector
	velocity Vector
	destroy  bool
	end      time.Time
}

// Draw ...
func (b *Baloon) Draw(ctx js.Value) {
	target := baloon
	if b.destroy {
		target = spark
	}
	x := int(b.pos.X - target.Get("width").Float()/2)
	y := int(b.pos.Y - target.Get("height").Float()/2)
	ctx.Call("drawImage", target, x, y)
}

// Destroy ...
func (b *Baloon) Destroy() bool {
	if !b.destroy {
		b.destroy = true
		b.end = time.Now()
		return true
	}
	return false
}

// Move ...
func (b *Baloon) Move(dt float64) {
	if !b.end.IsZero() && time.Since(b.end) > time.Second {
		delete(Baloons, b)
		return
	}
	v := b.velocity
	if b.destroy {
		v = v.Mul(0.3)
	}
	b.pos.X += v.X * dt
	b.pos.Y += v.Y * dt
	if b.pos.Y < -300 {
		log.Println("del:", b)
		delete(Baloons, b)
	}
}

// HitTest ...
func (b *Baloon) HitTest(d Vector) bool {
	return d.Distance(b.pos) < 60.0
}

// NewBaloon ...
func NewBaloon(x, y int) *Baloon {
	b := &Baloon{
		pos:      Vector{X: float64(x), Y: float64(y)},
		velocity: Vector{X: 0.0, Y: -0.3},
	}
	log.Println("add:", b)
	Baloons[b] = true
	return b
}
