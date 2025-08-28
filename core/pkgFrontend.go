package main

import (
	"backend/services/clients"
	backendscopes "backend/services/scopes"
	"core/example"
	"core/ping"
	"core/tacticalmap"
	"core/triangle"
	"errors"
	"fmt"
	"frontend/engine/systems/mainpipeline"
	"frontend/engine/tools/worldtexture"
	frontendapi "frontend/services/api"
	frontendtcp "frontend/services/api/tcp"
	"frontend/services/assets"
	"frontend/services/backendconnection"
	"frontend/services/backendconnection/localconnector"
	"frontend/services/colliders"
	"frontend/services/colliders/shapes"
	"frontend/services/console"
	"frontend/services/dbpkg"
	"frontend/services/ecs"
	"frontend/services/frames"
	"frontend/services/media"
	windowapi "frontend/services/media/window"
	"frontend/services/scenes"
	frontendscopes "frontend/services/scopes"
	"os"
	"path/filepath"
	"shared/services/api"
	"shared/services/logger"
	"shared/services/uuid"
	"shared/utils/connection"

	"github.com/go-gl/gl/v4.5-core/gl"
	// "github.com/go-gl/glfw/v3.3/glfw"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

func frontendDic(
	backendC ioc.Dic,
	sharedPkg SharedPkg,
) ioc.Dic {
	// sdl + opengl

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(fmt.Errorf("Failed to initialize SDL: %s", err))
	}

	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 4 /* 3 */)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1 /* 3 */)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1) // Essential for GLSwap
	sdl.GLSetAttribute(sdl.GL_DEPTH_SIZE, 24)  // Good practice for depth testing
	window, err := sdl.CreateWindow(
		"texhec",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600,
		sdl.WINDOW_SHOWN|sdl.WINDOW_OPENGL,
	)
	if err != nil {
		panic(fmt.Errorf("Failed to create window: %s", err))
	}

	ctx, err := window.GLCreateContext()
	if err != nil {
		panic(fmt.Errorf("Failed to create gl context: %s", err))
	}
	if err := gl.Init(); err != nil {
		panic(fmt.Errorf("could not initialize OpenGL: %v", err))
	}
	gl.Enable(gl.DEPTH_TEST)
	gl.Disable(gl.CULL_FACE)
	gl.DepthFunc(gl.LESS)
	// gl.DepthFunc(gl.GREATER)
	if err := window.GLMakeCurrent(ctx); err != nil {
		panic(fmt.Errorf("could not make OpenGL context current: %v", err))
	}
	sdl.GLSetSwapInterval(0)

	// path

	engineDir, err := os.Getwd()
	if err != nil {
		panic(errors.Join(errors.New("current wordking direcotry"), err))
	}
	// parent of both /backend and /frontend directory
	engineDir = filepath.Dir(engineDir)
	storage := filepath.Join(engineDir, "user_storage", "frontend")

	pkgs := []ioc.Pkg{
		sharedPkg,
		api.Package(func(c ioc.Dic) ioc.Dic { return c }),
		assets.Package(),
		colliders.Package(),
		shapes.Package(),
		logger.Package(true, func(c ioc.Dic, message string) {
			ioc.Get[console.Console](c).PrintPermanent(message)
		}),
		dbpkg.Package(fmt.Sprintf("%s/db.sql", storage)),
		frontendtcp.Package("tcp"),
		frontendapi.Package(),
		localconnector.Package(func(clientCon connection.Connection) connection.Connection {
			backendC := backendC.Scope(backendscopes.UserSession)
			client := clients.NewClient(
				clients.ClientID(ioc.Get[uuid.Factory](backendC).NewUUID().String()),
				clientCon,
			)
			sClient := ioc.Get[clients.SessionClient](backendC)
			sClient.UseClient(client)

			return ioc.Get[connection.Connection](backendC)
		}),
		backendconnection.Package(func(c ioc.Dic) connection.Connection {
			return ioc.Get[localconnector.Connector](c).Connect()
		}),
		console.Package(),
		media.Package(
			windowapi.Package(window, ctx),
		),
		ecs.Package(),
		// frames.Package(60),
		frames.Package(10000),
		scenes.Package(),
		frontendscopes.Package(),

		mainpipeline.Package(),
		worldtexture.Package(),

		// mods
		ping.FrontendPackage(),
		tacticalmap.FrontendPackage(),
		triangle.FrontendPackage(),
		example.FrontendPackage(),
	}

	b := ioc.NewBuilder()
	for _, pkg := range pkgs {
		pkg.Register(b)
	}

	return b.Build()
}
