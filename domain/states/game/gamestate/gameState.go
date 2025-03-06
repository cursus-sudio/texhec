package gamestate

import (
	"domain/common/models"
	"domain/states/game/gamedefinitions"
	"domain/states/game/gamevo"
	"domain/states/kernel"
	"domain/states/kernel/vo"
	"os/user"

	"github.com/ogiusek/null"
)

type FogOfWar interface {
	models.ModelBase
	Setter() (FogOfWarSetter, error)

	NeverSeenTiles() []gamevo.Pos
	CoveredTiles() map[gamevo.Pos]Tile
	SeesTiles() map[gamevo.Pos]Tile
}

type Research interface {
	models.ModelBase
	Setter() (ResearchSetter, error)

	Definition() kernel.ResearchDefinition
	Finished() bool
	Progress() vo.Budget
}

type Player interface {
	models.ModelBase
	Setter() (PlayerSetter, error)

	User() *user.User

	FogOfWar() FogOfWar
	Brain() Building
	Buildings() []Building
	Battalions() []TroopBattalion

	Fraction() gamedefinitions.Fraction
	TakesTurn() bool

	Studies() []Research
	Study(models.ModelId) null.Nullable[Research]
}

type TroopBattalion interface {
	models.ModelBase
	Setter() (TroopBattalionSetter, error)

	Definition() gamedefinitions.Troop
	Pos() gamevo.Pos

	Owner() Player
	Budget() vo.TroopBudget
	BrigadeSize() uint

	Skills() []gamedefinitions.TroopSkill
}

type Building interface {
	Setter() (Building, error)

	Pos() gamevo.Pos
	Definition() gamedefinitions.Resource
	Hp() vo.Hp
}

type Resource interface {
	Setter() (ResourceSetter, error)

	Pos() gamevo.Pos
	Definition() gamedefinitions.Resource
}

type Tile interface {
	Setter() (TileSetter, error)

	Pos() gamevo.Pos
	Definition() gamedefinitions.Tile
	Resource() null.Nullable[Resource]
	Building() null.Nullable[Building]
	Battalions() []TroopBattalion
	AttackerBattalions() []TroopBattalion
}

type Map interface {
	Setter() (MapSetter, error)

	Size() gamevo.Size
	GetTile(gamevo.Pos) Tile
}

// SERVICE
// TODO move to other file
type SizePositions interface {
	AllPositions(gamevo.Size) []gamevo.Pos
}

type GameState interface {
	models.ModelBase
	Setter() (GameStateSetter, error)

	Map() Map

	Players() []Player
	GetPlayer(models.ModelId) null.Nullable[Player]

	Battalions() []TroopBattalion
	GetBattaion(models.ModelId) null.Nullable[TroopBattalion]

	Buildings() []Building
	GetBuilding(models.ModelId) null.Nullable[Building]

	Resources() []Resource
	GetResource(models.ModelId) null.Nullable[Resource]
}
