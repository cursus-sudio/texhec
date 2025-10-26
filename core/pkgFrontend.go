package main

import (
	"backend/services/clients"
	backendscopes "backend/services/scopes"
	gameassets "core/assets"
	gamescenes "core/scenes"
	menuscene "core/scenes/menu"
	"core/src/ping"
	"core/src/tacticalmap"
	"core/src/tile"
	"errors"
	"fmt"
	"frontend/engine/components/groups"
	"frontend/engine/components/projection"
	"frontend/engine/components/text"
	"frontend/engine/systems/genericrenderer"
	"frontend/engine/systems/text"
	"frontend/engine/tools/broadcollision"
	"frontend/engine/tools/cameras"
	frontendapi "frontend/services/api"
	frontendtcp "frontend/services/api/tcp"
	"frontend/services/assets"
	"frontend/services/backendconnection"
	"frontend/services/backendconnection/localconnector"
	"frontend/services/console"
	"frontend/services/dbpkg"
	"frontend/services/frames"
	"frontend/services/graphics/texturearray"
	"frontend/services/media"
	windowapi "frontend/services/media/window"
	"frontend/services/scenes"
	frontendscopes "frontend/services/scopes"
	"os"
	"path/filepath"
	"shared/services/api"
	"shared/services/datastructures"
	"shared/services/logger"
	"shared/services/uuid"
	"shared/utils/connection"

	"github.com/go-gl/gl/v4.5-core/gl"
	"golang.org/x/image/font/opentype"

	// "github.com/go-gl/glfw/v3.3/glfw"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

func frontendDic(
	backendC ioc.Dic,
	sharedPkg SharedPkg,
) ioc.Dic {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(fmt.Errorf("Failed to initialize SDL: %s", err))
	}

	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 4 /* 3 */)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1 /* 3 */)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1) // Essential for GLSwap
	sdl.GLSetAttribute(sdl.GL_DEPTH_SIZE, 24)  // Good practice for depth testing

	// sdl.GLSetAttribute(sdl.GL_MULTISAMPLEBUFFERS, 1)
	// sdl.GLSetAttribute(sdl.GL_MULTISAMPLESAMPLES, 4)
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

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

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
		// ecs.Package(), // scenes register world so ecs package isn't registered
		frames.Package(60),
		// frames.Package(10000),
		scenes.Package(),
		frontendscopes.Package(),
		cameras.Package(),
		projection.Package(),

		genericrenderersys.Package(),
		broadcollision.Package(),

		texturearray.Package(),
		tile.Package(100, -1., groups.DefaultGroups()),

		textsys.Package(
			text.FontFamily{FontFamily: gameassets.FontAssetID},
			text.FontSize{FontSize: 16},
			// text.Overflow{Visible: false},
			text.Break{Break: text.BreakWord},
			text.TextAlign{Vertical: 0, Horizontal: 0},
			func() datastructures.SparseSet[rune] {
				set := datastructures.NewSparseSet[rune]()
				for i := int32('a'); i <= int32('z'); i++ {
					set.Add(rune(i))
				}
				for i := int32('A'); i <= int32('Z'); i++ {
					set.Add(rune(i))
				}
				for i := int32('0'); i <= int32('9'); i++ {
					set.Add(rune(i))
				}
				for i := int32('!'); i <= int32('/'); i++ {
					set.Add(rune(i))
				}
				for i := int32(':'); i <= int32('@'); i++ {
					set.Add(rune(i))
				}
				for i := int32('['); i <= int32('`'); i++ {
					set.Add(rune(i))
				}
				for i := int32('{'); i <= int32('~'); i++ {
					set.Add(rune(i))
				}
				set.Add(' ')

				return set
			}(),
			opentype.FaceOptions{
				Size: 64,
				DPI:  72,
			},
			52,
		),

		// mods
		ping.FrontendPackage(),
		tacticalmap.FrontendPackage(),

		gamescenes.Package(),
		gameassets.Package(),
		menuscene.Package(),
	}

	b := ioc.NewBuilder()
	for _, pkg := range pkgs {
		pkg.Register(b)
	}

	return b.Build()
}
