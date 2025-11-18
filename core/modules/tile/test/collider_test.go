package test

import (
	"core/modules/tile"
	"core/modules/tile/internal/tilecollider"
	"testing"
)

func TestCollider(t *testing.T) {
	const (
		L1 tile.Layer = iota
		L2
	)

	collider := tilecollider.NewCollider()
	if collider.Has(L1) {
		t.Error("empty collider already has l1. read or ctor fault")
		return
	}
	if collider.Has(L2) {
		t.Error("empty collider already has l2. read or ctor fault")
		return
	}
	collider.Add(L1)
	if !collider.Has(L1) {
		t.Error("collider should have l1. add fault")
		return
	}
	if collider.Has(L2) {
		t.Error("collider shouldn't have l2. add fault")
		return
	}
	collider.Add(L2)
	if !collider.Has(L1) {
		t.Error("collider should have l1. add fault")
		return
	}
	if !collider.Has(L2) {
		t.Error("collider should have l2. add fault")
		return
	}
	collider.Remove(L1)
	if collider.Has(L1) {
		t.Error("collider shouldn't have l1. remove fault")
		return
	}
	if !collider.Has(L2) {
		t.Error("collider should have l2. remove fault")
		return
	}
	collider.Remove(L2)
	if collider.Has(L1) {
		t.Error("collider shouldn't have l1. remove fault")
		return
	}
	if collider.Has(L2) {
		t.Error("collider shouldn't have l2. remove fault")
		return
	}
}
