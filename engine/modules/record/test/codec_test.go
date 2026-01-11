package test

import (
	"engine/modules/record"
	"maps"
	"slices"
	"testing"
)

func TestCodec(t *testing.T) {
	s := NewSetup()
	initialState := Component{Counter: 6}

	entity := s.World.NewEntity()
	s.ComponentArray.Set(entity, initialState)

	originalRecording := s.Record.UUID().GetState(s.Config)

	encodedRecording, err := s.Codec.Encode(originalRecording)
	if err != nil {
		t.Error(err)
		return
	}

	decodedRecording, err := s.Codec.Decode(encodedRecording)
	if err != nil {
		t.Error(err)
		return
	}

	comparedRecording, ok := decodedRecording.(record.UUIDRecording)
	if !ok {
		t.Error("decoded recording doesn't have encoded recording type")
		return
	}

	// UUIDEntities
	if !maps.EqualFunc(originalRecording.Entities, comparedRecording.Entities, func(v1, v2 []any) bool {
		return slices.Equal(v1, v2)
	}) {
		t.Errorf(
			"Entities don't match expected %v has %v",
			originalRecording.Entities,
			comparedRecording.Entities,
		)
	}
}
