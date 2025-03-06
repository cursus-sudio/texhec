package battle

import (
	"domain/common/models"
)

type BattleState struct {
	models.ModelBase
	Players []*Player
	Map     *Map
}
