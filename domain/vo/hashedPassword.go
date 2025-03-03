package vo

import "github.com/ogiusek/ioc"

type Hash []byte

func (hash *Hash) Valid(c ioc.Dic) []error {
	// TODO
	return nil
}

type IHasher interface {
	Hash([]byte) Hash
	Compare(Hash, Hash) (match bool)
}
