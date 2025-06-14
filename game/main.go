package main

import (
	backendsrc "backend/src"
	"backend/src/backendapi"
	"backend/src/modules"
	"backend/src/modules/tacticalmap"
	"backend/src/utils"
	"backend/src/utils/services"
	"backend/src/utils/services/scopecleanup"
	"frontend/src/gameloop"
	"log"

	"github.com/ogiusek/ioc"
)

var pkgs []ioc.Pkg = []ioc.Pkg{
	backendsrc.Package(
		utils.Package(
			services.Package(
				scopecleanup.Package(),
			),
		),
		modules.Package(
			tacticalmap.Package(),
		),
		backendapi.Package(),
	),
}

func main() {
	c := ioc.NewContainer()

	for _, pkg := range pkgs {
		pkg.Register(c)
	}

	backend := ioc.Get[backendapi.Backend](c)
	tacticalMap := backend.TacticalMap()

	if err := tacticalMap.Create(tacticalmap.CreateArgs{
		Tiles: []tacticalmap.Tile{
			{Pos: tacticalmap.Pos{X: 1, Y: 1}},
			{Pos: tacticalmap.Pos{X: 1, Y: 2}},
		},
	}); err != nil {
		log.Panicf("p1 %s", err.Error())
	}

	tiles, err := tacticalMap.GetMap()
	if err != nil {
		log.Panicf("p2 %s", err.Error())
	}

	log.Printf("%v", tiles)

	gameloop.StartGameLoop(c)
}
