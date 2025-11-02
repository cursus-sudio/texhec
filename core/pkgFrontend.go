package main

import (
	"backend/services/clients"
	backendscopes "backend/services/scopes"
	gameassets "core/assets"
	"core/modules/fpslogger/pkg"
	tilepkg "core/modules/tile/pkg"
	gamescenes "core/scenes"
	creditsscene "core/scenes/credits"
	gamescene "core/scenes/game"
	menuscene "core/scenes/menu"
	settingsscene "core/scenes/settings"
	"errors"
	"fmt"
	"frontend/modules/anchor/pkg"
	"frontend/modules/camera/pkg"
	"frontend/modules/collider/pkg"
	"frontend/modules/genericrenderer/pkg"
	"frontend/modules/groups"
	"frontend/modules/groups/pkg"
	"frontend/modules/inputs/pkg"
	"frontend/modules/render/pkg"
	"frontend/modules/scenes/pkg"
	"frontend/modules/text"
	"frontend/modules/text/pkg"
	"frontend/modules/transform/pkg"
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
		media.Package(window, ctx),
		// ecs.Package(), // scenes register world so ecs package isn't registered
		frames.Package(60),
		// frames.Package(10000),
		scenes.Package(),
		frontendscopes.Package(),

		texturearray.Package(),
		tilepkg.Package(100, -1., groups.EmptyGroups().Ptr().Enable(gamescene.GameGroup).Val()),

		//

		// engine packages
		anchorpkg.Package(),
		camerapkg.Package(),
		colliderpkg.Package(),
		genericrendererpkg.Package(),
		groupspkg.Package(),
		inputspkg.Package(),
		renderpkg.Package(),
		scenespkg.Package(),
		textpkg.Package(
			text.FontFamilyComponent{FontFamily: gameassets.FontAssetID},
			text.FontSizeComponent{FontSize: 16},
			// text.Overflow{Visible: false},
			text.BreakComponent{Break: text.BreakWord},
			text.TextAlignComponent{Vertical: 0, Horizontal: 0},
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
		transformpkg.Package(),

		// game packages
		fpsloggerpkg.Package(),

		gamescenes.Package(),
		gameassets.Package(),

		creditsscene.Package(),
		gamescene.Package(),
		menuscene.Package(),
		settingsscene.Package(),
	}

	b := ioc.NewBuilder()
	for _, pkg := range pkgs {
		pkg.Register(b)
	}

	return b.Build()
}
