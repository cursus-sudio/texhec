package vo

import (
	"encoding/json"
)

type Hp struct {
	hp uint
}

func NoHP() Hp {
	return Hp{hp: uint(0)}
}

func NewHP(health uint) Hp {
	return Hp{hp: health}
}

func (hp Hp) IsPositive() bool {
	return hp.hp != 0 // hp.health > 0
}

func (hp Hp) IsGreatherThan(than Hp) bool {
	return hp.hp > than.hp
}

func (hp Hp) Sum(add Hp) Hp {
	hp.hp += add.hp
	return hp
}

func (hp Hp) Deal(dmg Hp) (hpLeft Hp, overflow Hp) {
	if hp.hp > dmg.hp {
		return Hp{hp: hp.hp - dmg.hp}, NoHP()
	} else {
		return NoHP(), Hp{hp: dmg.hp - hp.hp}
	}
}

func (hp Hp) Multiply(multiplier float32) Hp {
	hp.hp = uint(multiplier * float32(hp.hp))
	return hp
}

func (o *Hp) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.hp)
}

func (o *Hp) UnmarshalJSON(data []byte) error {
	var hp uint
	if err := json.Unmarshal(data, &hp); err != nil {
		return err
	}
	o.hp = hp
	return nil
}
