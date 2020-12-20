package main

import "time"

// State ...
type State struct {
	Time          time.Time
	State         int
	Valid         bool
	Distance      float64
	Score         int
	LeftWrist     Vector
	RightWrist    Vector
	Nose          Vector
	LeftShoulder  Vector
	RightShoulder Vector
}
