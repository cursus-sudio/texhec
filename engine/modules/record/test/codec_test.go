package test

import (
	"engine/modules/record"
	"maps"
	"slices"
	"testing"
)

func TestCodec(t *testing.T) {
	world := NewSetup()
	initialState := Component{Counter: 6}

	entity := world.NewEntity()
	world.ComponentArray.Set(entity, initialState)

	originalRecording := world.Record().UUID().GetState(world.Config)

	encodedRecording, err := world.Codec.Encode(originalRecording)
	if err != nil {
		t.Error(err)
		return
	}

	decodedRecording, err := world.Codec.Decode(encodedRecording)
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
	if !maps.Equal(originalRecording.UUIDEntities, comparedRecording.UUIDEntities) {
		t.Errorf(
			"UUIDEntities mismatch. wanted %v has %v",
			originalRecording.UUIDEntities,
			comparedRecording.UUIDEntities,
		)
	}

	// EntitiesUUIDs
	// if !slices.Equal(originalRecording.EntitiesUUIDs.GetValues(), comparedRecording.EntitiesUUIDs.GetValues()) {
	// 	t.Errorf(
	// 		"EntitiesUUIDs values mismatch. %v != %v",
	// 		originalRecording.EntitiesUUIDs.GetValues(),
	// 		comparedRecording.EntitiesUUIDs.GetValues(),
	// 	)
	// }

	// RemovedEntities
	if !slices.Equal(originalRecording.RemovedEntities.GetIndices(), comparedRecording.RemovedEntities.GetIndices()) {
		t.Errorf(
			"RemovedEntities mismatch. %v != %v",
			originalRecording.RemovedEntities.GetIndices(),
			comparedRecording.RemovedEntities.GetIndices(),
		)
	}

	// Arrays
	if !maps.EqualFunc(originalRecording.Arrays, comparedRecording.Arrays, func(v1, v2 record.ArrayRecording) bool {
		return slices.Equal(v1.GetIndices(), v2.GetIndices()) &&
			slices.Equal(v1.GetValues(), v2.GetValues())
	}) {
		t.Errorf(
			"Arrays mismatch. %v != %v",
			originalRecording.Arrays,
			comparedRecording.Arrays,
		)
	}

}
