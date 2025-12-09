package main

import (
	_ "embed"
	"engine/services/ecs"
	appruntime "engine/services/runtime"
	"runtime"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/ioc/v2"
)

func main() {
	print("started\n")

	runtime.LockOSThread()

	sharedPkg := SharedPackage()

	c := frontendDic(
		sharedPkg,
	)

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	// load world before starting timer
	ioc.Get[ecs.World](c)
	frontendRuntime := ioc.Get[appruntime.Runtime](c)
	frontendRuntime.Run()
}
