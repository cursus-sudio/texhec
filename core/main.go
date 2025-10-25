package main

import (
	"core/src/ping"
	"core/src/tacticalmap"
	"fmt"
	frontendtcp "frontend/services/api/tcp"
	"frontend/services/backendconnection"
	"os"
	"runtime"
	appruntime "shared/services/runtime"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/relay/v2"
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
		// go func() {
		// 	time.Sleep(time.Second / 10)
		// 	backendC := backendC.Scope(scopes.Request)
		// 	s := ioc.Get[saves.Saves](backendC)
		// 	factory := ioc.Get[saves.SaveMetaFactory](backendC)
		// 	err := s.NewSave(factory.New("very funny save\n"))
		// 	ioc.Get[scopes.RequestService](backendC).Clean(scopes.NewRequestEndArgs(err))
		// 	ioc.Get[logger.Logger](backendC).Info("saved")
		// }()
		backendRuntime.Run()
		return
	}

	c := frontendDic(
		backendC,
		sharedPkg,
	)

	if false { // connect
		tcpConnect := ioc.Get[frontendtcp.Connect](c)
		err := tcpConnect.Connect("localhost:8080")
		if err != nil {
			panic(err)
		}
	}
	{ // pinging backend
		backend := ioc.Get[backendconnection.Backend](c).Connection()
		r := backend.Relay()
		res, err := relay.Handle(r, ping.PingReq{ID: 2077})
		fmt.Printf("client recieved ping res is %v\nerr is %s\n", res, err)
	}
	{
		r := ioc.Get[backendconnection.Backend](c).Connection().Relay()
		res, err := relay.Handle(r, tacticalmap.NewCreateReq(
			tacticalmap.CreateArgs{
				Tiles: []tacticalmap.Tile{
					{Pos: tacticalmap.Pos{X: 7, Y: 13}},
				},
			},
		))
		fmt.Printf("create res is %v\nerr is %s\n", res, err)
	}
	{
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	}

	frontendRuntime := ioc.Get[appruntime.Runtime](c)
	// go func() {
	// 	time.Sleep(time.Second / 10)
	// 	frontendRuntime.Stop()
	// }()
	frontendRuntime.Run()
}
