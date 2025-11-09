package main

import (
	_ "embed"
	"os"
	"runtime"
	"shared/services/ecs"
	appruntime "shared/services/runtime"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/ioc/v2"
)

func main() {
	print("started\n")

	runtime.LockOSThread()

	isServer := false
	for _, arg := range os.Args {
		if arg == "server" {
			isServer = true
			break
		}
	}

	sharedPkg := SharedPackage()

	backendC := backendDic(sharedPkg)
	if isServer {
		backendRuntime := ioc.Get[appruntime.Runtime](backendC)
		backendRuntime.Run()
		return
	}

	c := frontendDic(
		backendC,
		sharedPkg,
	)

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	// load world before starting timer
	ioc.Get[ecs.World](c)
	frontendRuntime := ioc.Get[appruntime.Runtime](c)
	frontendRuntime.Run()
}
