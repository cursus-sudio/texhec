package main

import (
	gamescenes "core/scenes"
	_ "embed"
	"engine/modules/scene"
	appruntime "engine/services/runtime"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

// golangci-lint run --fix
func main() {
	print("started\n")

	{ // go tool pprof -http=:8080 cpu.pprof.cp
		name := ""
		if len(os.Args) > 1 {
			name = os.Args[1]
		}
		f, err := os.Create(fmt.Sprintf("cpu.pprof%v.cp", name))
		if err != nil {
			panic(err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}

	runtime.LockOSThread()

	c := getDic()

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)

	// load world before starting timer
	events.Emit(ioc.Get[events.Events](c), scene.NewChangeSceneEvent(gamescenes.GameID))
	frontendRuntime := ioc.Get[appruntime.Runtime](c)
	frontendRuntime.Run()
}
