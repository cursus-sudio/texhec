package registry

import (
	"engine/services/ecs"
	"fmt"
)

var (
	ErrExpectedPointerToAStruct error = fmt.Errorf("expected pointer to a struct")
)

type Service interface {
	Register(structTagKey string, handler func(structTagValue string) ecs.EntityID)

	// can return ErrExpectedPointerToAStruct
	Populate(any) error
}

func GetRegistry[Registry any](s Service) (Registry, error) {
	var r Registry
	err := s.Populate(&r)
	return r, err
}
