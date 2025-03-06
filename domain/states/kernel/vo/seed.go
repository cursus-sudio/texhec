package vo

import (
	"crypto/rand"

	"github.com/ogiusek/ioc"
)

type Seed []byte

func NewSeed(c ioc.Dic) Seed {
	seed := make([]byte, 16)
	rand.Read(seed)
	return Seed(seed)
}
