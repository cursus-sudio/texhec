package vo

import (
	"encoding/json"
)

type HP struct {
	hp uint
}

func NoHP() HP {
	return HP{hp: uint(0)}
}
func NewHP(health uint) HP {
	return HP{hp: health}
}

func (hp HP) IsPositive() bool {
	return hp.hp != 0 // hp.health > 0
}

func (hp HP) IsGreatherThan(than HP) bool {
	return hp.hp > than.hp
}

func (hp HP) Sum(add HP) HP {
	hp.hp += add.hp
	return hp
}

func (hp HP) Deal(dmg HP) (hpLeft HP, overflow HP) {
	if hp.hp > dmg.hp {
		return HP{hp: hp.hp - dmg.hp}, NoHP()
	} else {
		return NoHP(), HP{hp: dmg.hp - hp.hp}
	}
}

func (hp HP) Multiply(multiplier float32) HP {
	hp.hp = uint(multiplier * float32(hp.hp))
	return hp
}

func (o *HP) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.hp)
}

func (o *HP) UnmarshalJSON(data []byte) error {
	var hp uint
	if err := json.Unmarshal(data, &hp); err != nil {
		return err
	}
	o.hp = hp
	return nil
}
