package program

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-gl/gl/v4.5-core/gl"
)

var (
	ErrNotALocation    error = errors.New("expected 'int32' for location")
	ErrInvalidLocation error = errors.New("invalid location")
)

type Parameter struct {
	Name  uint32
	Value int32
}

func compileProgram(program uint32, parameters []Parameter) error {
	for _, p := range parameters {
		gl.ProgramParameteri(program, p.Name, p.Value)
	}

	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		infoLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(infoLog))

		return fmt.Errorf("failed to link program: %v", infoLog)
	}

	return nil
}

// fields with tag `uniform:"uniformName"` will have automatically generated uniform
func createLocations(t reflect.Type, program uint32) (any, error) {
	val := reflect.New(t).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)
		uniformName := field.Tag.Get("uniform")

		if uniformName == "" {
			continue
		}

		if field.Type.Kind() != reflect.Int32 {
			err := errors.Join(
				ErrNotALocation,
				fmt.Errorf(
					"field \"%s.%s\" isn't int32",
					typ.String(),
					field.Name,
				),
			)
			return nil, err
		}

		location := gl.GetUniformLocation(program, gl.Str(uniformName+"\x00"))
		if location == -1 {
			err := errors.Join(
				ErrInvalidLocation,
				fmt.Errorf("uniform \"%s\" doesn't exist in shader program", uniformName),
			)
			return nil, err
		}
		fieldValue.SetInt(int64(location))
	}
	return val.Interface(), nil
}
