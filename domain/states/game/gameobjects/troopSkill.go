package gameobjects

import (
	"domain/common/events"
	"domain/states/game/gamedefinitions"
	"log"

	"github.com/ogiusek/ioc"
)

type TroopSkillPayload struct {
	c     ioc.Dic
	game  *GameState
	troop *Troop
}

func (skill *TroopSkillPayload) Game() *GameState {
	return skill.game
}

func (skill *TroopSkillPayload) Troop() *Troop {
	return skill.troop
}

func TroopSuicideSkill(skill TroopSkillPayload) {
	skill.troop.SetAmount(skill.c, 0)
}

func RegisterTroopSuicideSkill(c ioc.Dic) {
	const topic = "xxxxx-xxxxx-xxxx-xxxx"
	e := ioc.Get[events.Events[TroopSkillPayload]](c)
	e.RegisterHandler(topic, func(e events.Event[TroopSkillPayload]) { TroopSuicideSkill(e.Payload()) })
	// var troopSkill TroopSkill
	// event := events.NewEvent(troopSkill)
	// make this use types
	// manager.Listen(event.Topic, func(a any) { TroopSuicideSkill(c, a.(TroopSkill)) })
}

type exRepo interface {
	GetSkills() []gamedefinitions.TroopSkill
}

func OnStart(c ioc.Dic) {
	skills := ioc.Get[exRepo](c)
	events := ioc.Get[events.Events[TroopSkillPayload]](c)
	for _, skill := range skills.GetSkills() {
		if !events.IsTaken(skill.SkillId()) {
			log.Panic("missing registered skill")
		}
	}
}
