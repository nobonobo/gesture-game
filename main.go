package main

import (
	"log"
	"syscall/js"

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

func main() {
	poseNet := posenet.New(posenet.DefaultSingleConfig)
	playing := false
	log.Println("loaded")
	var cb js.Func
	ctx := js.Null()
	video := js.Null()
	width := 1280
	height := 720
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
				poses, err := poseNet.EstimateSinglePose(nil)
				if err != nil {
					log.Println(err)
					return
				}
				_ = poses
				console.Call("log", poses)
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
