package program

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-gl/gl/v4.5-core/gl"
)

var (
	ErrProgramHasOtherLocations error = errors.New("invalid program locations type")
)

type Program interface {
	ID() uint32
	Locations(reflect.Type) (any, error)
	Bind()
	Release()
}

type program struct {
	id            uint32
	locationsType reflect.Type
	locations     any
}

func NewProgram(p uint32, parameters []Parameter) (Program, error) {
	err := compileProgram(p, parameters)
	if err != nil {
		return nil, err
	}
	return &program{id: p}, nil
}

func (p *program) ID() uint32 { return p.id }

func (p *program) Bind() {
	gl.UseProgram(p.id)
}

func (p *program) Locations(t reflect.Type) (any, error) {
	if p.locations != nil {
		if p.locationsType != t {
			err := errors.Join(
				ErrProgramHasOtherLocations,
				fmt.Errorf("requested \"%s\" but program has \"%s\"", t.String(), p.locationsType.String()),
			)
			return nil, err
		}
		return p.locations, nil
	}
	locations, err := createLocations(t, p.id)
	if err != nil {
		return nil, err
	}
	p.locations = locations
	p.locationsType = t
	return locations, nil
}

func (p *program) Release() {
	gl.DeleteProgram(p.id)
}

func GetProgramLocations[Locations any](p Program) (Locations, error) {
	locations, err := p.Locations(reflect.TypeFor[Locations]())
	if err != nil {
		var l Locations
		return l, err
	}
	return locations.(Locations), nil
}
