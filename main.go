package main

import (
	"fmt"
	"log"
	"math/rand"
	"syscall/js"
	"time"

	"github.com/nobonobo/posenet"
	"github.com/nobonobo/spago"
	"github.com/nobonobo/spago/dispatcher"
	"github.com/nobonobo/spago/jsutil"
	"github.com/nobonobo/spago/router"

	"gesture-game/actions"
	"gesture-game/views"
)

var (
	document = js.Global().Get("document")
	console  = js.Global().Get("console")
	white    = "rgb(255,255,255)"
)

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

	actx := js.Global().Get("AudioContext").New()
	resp, err := jsutil.Fetch("assets/clap.mp3", nil)
	if err != nil {
		log.Fatal(err)
	}
	buff, err := jsutil.Await(resp.Call("arrayBuffer"))
	if err != nil {
		log.Fatal(err)
	}
	clapSound, err := jsutil.Await(actx.Call("decodeAudioData", buff))
	if err != nil {
		log.Fatal(err)
	}
	console.Call("log", clapSound)

	poseNet := posenet.New(conf)
	playing := false
	log.Println("loaded")
	var cb js.Func
	ctx := js.Null()
	video := js.Null()
	width := 640
	height := 480
	lastTick := float64(0.0)
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		current := args[0].Float()
		dt := float64(0.0)
		if lastTick != 0 {
			dt = current - lastTick
		}
		lastTick = current
		go func() {
			terminate := false
			if playing {
				if ctx.IsNull() {
					baloon = document.Call("getElementById", "baloon")
					spark = document.Call("getElementById", "spark")
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
				drawText(ctx, fmt.Sprintf("てんすう: %d", state.Score), 0, 24, 24, white)
				remain := 30 - int(time.Since(state.Time)/time.Second)
				if remain < -3 {
					terminate = true
				}
				if remain < 0 {
					drawTextCenter(ctx, "しゅうりょう！", 320, 200, 32, "rgb(100,100,255)")
					remain = 0
				}
				drawTextRight(ctx, fmt.Sprintf("のこり: %ds", remain), 640, 20, 20, white)
				if state.State == 1 {
					if rand.Int()%(remain+3) < 1 {
						dispatcher.Dispatch(actions.Spawn)
					}
					ctx.Call("save")
					for b := range Baloons {
						b.Draw(ctx)
						b.Move(dt)
						if remain > 0 {
							if b.HitTest(state.LeftWrist) ||
								b.HitTest(state.RightWrist) {
								if b.Destroy() {
									dispatcher.Dispatch(actions.Hit)
									state.Score++
								}
							}
						}
					}
					ctx.Call("restore")
				}
				if state.Distance < 5 {
					drawTextCenter(ctx, "ちかい！ちかい！", 320, 200, 32, "rgb(100,100,255)")
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
								drawTextCenter(ctx, ts, 320, 240, 48, "rgb(200,200,255)")
							} else {
								state.State = 1
								state.Time = time.Now()
							}
						}
						drawDot(ctx, state.LeftWrist, "#ff0000")
						drawDot(ctx, state.RightWrist, "#00ff00")
					}
				}
				if terminate {
					dispatcher.Dispatch(actions.GameEnd)
				}
				poses, err := poseNet.EstimateSinglePose(nil)
				if err != nil {
					log.Println(err)
					return
				}
				valid := 0
				if poses.Get("score").Float() > 0.4 {
					keypoints := poses.Get("keypoints")
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
	dispatcher.Register(actions.Spawn, func(args ...interface{}) {
		x := rand.Int() % 640
		NewBaloon(x, 480)
	})
	js.Global().Set("spawn", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		dispatcher.Dispatch(actions.Spawn)
		return nil
	}))

	dispatcher.Register(actions.Hit, func(args ...interface{}) {
		source := actx.Call("createBufferSource")
		source.Set("buffer", clapSound)
		source.Call("connect", actx.Get("destination"))
		source.Call("start")
	})
	dispatcher.Register(actions.Refresh, func(args ...interface{}) {
		if c, ok := args[0].(spago.Component); ok {
			spago.Rerender(c)
		}
	})
	dispatcher.Register(actions.GameStart, func(args ...interface{}) {
		dispatcher.Dispatch(actions.Hit)
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
		spago.RenderBody(&views.Result{Score: state.Score})
	})
	r.Start()
	select {}
}
