package main

import (
	_ "embed"
	"engine/services/ecs"
	appruntime "engine/services/runtime"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/ioc/v2"
)

func Invalid() error { return nil }

// golangci-lint run --fix
func main() {
	Invalid()
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
	ioc.Get[ecs.World](c)
	frontendRuntime := ioc.Get[appruntime.Runtime](c)
	frontendRuntime.Run()
}
