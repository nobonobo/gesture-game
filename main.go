package main

import (
	"fmt"
	"log"
	"math"
	"syscall/js"
	"time"

	"github.com/nobonobo/posenet"
	"github.com/nobonobo/spago"
	"github.com/nobonobo/spago/dispatcher"
	"github.com/nobonobo/spago/router"

	"gesture-game/actions"
	"gesture-game/views"
)

var (
	document = js.Global().Get("document")
	console  = js.Global().Get("console")
)

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

func drawLine(ctx js.Value, begin, end Vector) {
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

func main() {
	conf := posenet.Config{
		Algorithm:       "single-pose",
		Architecture:    "MobileNetV1",
		OutputStride:    16,
		InputResolution: 300,
		Multiplier:      0.5,
		QuantBytes:      2,
	}
	state := State{}
	poseNet := posenet.New(conf)
	playing := false
	log.Println("loaded")
	var cb js.Func
	ctx := js.Null()
	video := js.Null()
	width := 640
	height := 480
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			if playing {
				if ctx.IsNull() {
					canvas := document.Call("getElementById", "output")
					canvas.Set("width", 640)
					canvas.Set("height", 480)
					width = canvas.Get("clientWidth").Int()
					height = canvas.Get("clientHeight").Int()
					log.Println("width:", width, "height:", height)
					ctx = canvas.Call("getContext", "2d")
					stream := document.Call("getElementById", "smallVideo").Get("srcObject")
					video = document.Call("createElement", "video")
					video.Set("id", "video")
					video.Get("style").Set("display", "none")
					video.Set("autoplay", true)
					video.Set("muted", true)
					video.Set("playsinline", true)
					video.Set("srcObject", stream)
					video.Set("width", 640)
					video.Set("height", 480)
					document.Get("body").Call("appendChild", video)
					video.Call("play")
				}
				ctx.Call("save")
				ctx.Call("scale", -1, 1)
				ctx.Call("translate", -640, 0)
				ctx.Call("drawImage", video, 0, 0)
				ctx.Call("restore")
				if state.Distance < 5 {
					drawText(ctx, "ちかい！ちかい！", 320, 200, 32, "rgb(100,100,255)")
					if state.State == 0 {
						state.Time = time.Now()
					}
				} else {
					if state.Valid {
						if state.State == 0 {
							tm := int(time.Since(state.Time))
							if tm < 4*int(time.Second) {
								t := 3 - tm/int(time.Second)
								ts := fmt.Sprintf("%d", t)
								if t == 0 {
									ts = "START!"
								}
								drawText(ctx, ts, 320, 240, 48, "rgb(200,200,255)")
							} else {
								state.State = 1
								state.Time = time.Now()
							}
						}
						drawDot(ctx, state.LeftWrist, "#ff0000")
						drawDot(ctx, state.RightWrist, "#00ff00")
					}
				}
				poses, err := poseNet.EstimateSinglePose(nil)
				if err != nil {
					log.Println(err)
					return
				}
				valid := 0
				if poses.Get("score").Float() > 0.4 {
					keypoints := poses.Get("keypoints")
					console.Call("log", keypoints)
					for i := 0; i < keypoints.Length(); i++ {
						part := keypoints.Index(i)
						if part.Get("score").Float() > 0.4 {
							switch part.Get("part").String() {
							case "nose":
								valid++
								state.Nose = NewVector(part.Get("position"))
							case "leftShoulder":
								valid++
								state.LeftShoulder = NewVector(part.Get("position"))
							case "rightShoulder":
								valid++
								state.RightShoulder = NewVector(part.Get("position"))
							case "leftWrist":
								valid++
								state.LeftWrist = NewVector(part.Get("position"))
							case "rightWrist":
								valid++
								state.RightWrist = NewVector(part.Get("position"))
							}
						}
					}
				}
				state.Valid = valid >= 5
				if state.Valid {
					state.Distance = 1000.0 / (state.Nose.Distance(state.LeftShoulder) +
						state.Nose.Distance(state.RightShoulder))
				}
			}
			js.Global().Call("requestAnimationFrame", cb)
		}()
		return nil
	})
	js.Global().Call("requestAnimationFrame", cb)
	dispatcher.Register(actions.Refresh, func(args ...interface{}) {
		if c, ok := args[0].(spago.Component); ok {
			spago.Rerender(c)
		}
	})
	dispatcher.Register(actions.GameStart, func(args ...interface{}) {
		go func() {
			router.Navigate("/play")
			poseNet.Start("smallVideo")
			playing = true
			state.Time = time.Now()
			state.State = 0
		}()
	})
	dispatcher.Register(actions.GameEnd, func(args ...interface{}) {
		go func() {
			playing = false
			ctx = js.Null()
			poseNet.Stop()
			router.Navigate("/result")
		}()
	})
	dispatcher.Register(actions.EstimatePose, func(args ...interface{}) {
		go func() {
			poses, err := poseNet.EstimateSinglePose(nil)
			if err != nil {
				log.Println(err)
				return
			}
			_ = poses
			//console.Call("log", poses)
		}()
	})
	log.Println(poseNet)
	r := router.New()
	r.Handle("/", func(key string) {
		log.Print("top")
		spago.RenderBody(&views.Top{})
	})
	r.Handle("/play", func(key string) {
		log.Print("playing")
		spago.RenderBody(&views.Playing{})
	})
	r.Handle("/result", func(key string) {
		log.Print("result")
		spago.RenderBody(&views.Result{})
	})
	r.Start()
	select {}
}
