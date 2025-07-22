package program

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-gl/gl/v4.5-core/gl"
)

func createProgram(vertexShader, fragmentShader uint32) (uint32, error) {
	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		infoLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(infoLog))

		return 0, fmt.Errorf("failed to link program: %v", infoLog)
	}

	return program, nil
}

// fields with tag `uniform:"uniformName"` will have automatically generated uniform
func createLocations[Locations any](program uint32) Locations {
	var l Locations

	val := reflect.ValueOf(&l).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)
		uniformName := field.Tag.Get("uniform")

		if field.Type.Kind() != reflect.Int32 || uniformName == "" {
			continue
		}

		// ignores is location -1
		location := gl.GetUniformLocation(program, gl.Str(uniformName+"\x00"))
		fieldValue.SetInt(int64(location))
	}
	return l
}
