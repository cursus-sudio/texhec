package test

import (
	"testing"
)

func BenchmarkEntityCodecEncode(b *testing.B) {
	world := NewSetup()

	entity := world.NewEntity()
	world.ComponentArray.Set(entity, Component{Counter: 6})

	originalRecording := world.Record().Entity().GetState(world.Config)

	for i := 0; i < b.N; i++ {
		_, _ = world.Codec.Encode(originalRecording)
	}
}
func BenchmarkEntityCodecDecode(b *testing.B) {
	world := NewSetup()

	entity := world.NewEntity()
	world.ComponentArray.Set(entity, Component{Counter: 6})

	originalRecording := world.Record().Entity().GetState(world.Config)

	encodedRecording, err := world.Codec.Encode(originalRecording)
	if err != nil {
		b.Error(err)
		return
	}

	for i := 0; i < b.N; i++ {
		_, _ = world.Codec.Decode(encodedRecording)
	}
}
func BenchmarkUUIDCodecEncode(b *testing.B) {
	world := NewSetup()

	entity := world.NewEntity()
	world.ComponentArray.Set(entity, Component{Counter: 6})

	originalRecording := world.Record().UUID().GetState(world.Config)

	for i := 0; i < b.N; i++ {
		_, _ = world.Codec.Encode(originalRecording)
	}
}
func BenchmarkUUIDCodecDecode(b *testing.B) {
	world := NewSetup()

	entity := world.NewEntity()
	world.ComponentArray.Set(entity, Component{Counter: 6})

	originalRecording := world.Record().UUID().GetState(world.Config)

	encodedRecording, err := world.Codec.Encode(originalRecording)
	if err != nil {
		b.Error(err)
		return
	}

	for i := 0; i < b.N; i++ {
		_, _ = world.Codec.Decode(encodedRecording)
	}
}
