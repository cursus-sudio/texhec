package transformtool

import (
	"engine/modules/transform"
	"engine/services/ecs"
)

func (t tool) Init() {
	onPosUpsert := func(ei []ecs.EntityID) {
		posTransaction := t.posArray.Transaction()
		for _, entity := range ei {
			parentObj := t.hierarchyTransaction.GetObject(entity)
			for _, child := range parentObj.Children().GetIndices() {
				mask, err := t.parentMaskArray.GetComponent(child)
				if err != nil {
					continue
				}
				if mask.RelativeMask&transform.RelativePos != 0 {
					posTransaction.TriggerChangeListener(child)
				}
			}
		}
		t.logger.Warn(ecs.FlushMany(posTransaction))
	}

	onRotUpsert := func(ei []ecs.EntityID) {
		posTransaction := t.posArray.Transaction()
		rotTransaction := t.posArray.Transaction()
		for _, entity := range ei {
			parentObj := t.hierarchyTransaction.GetObject(entity)
			for _, child := range parentObj.Children().GetIndices() {
				mask, err := t.parentMaskArray.GetComponent(child)
				if err != nil {
					continue
				}
				if mask.RelativeMask&transform.RelativePos != 0 {
					posTransaction.TriggerChangeListener(child)
				}
				if mask.RelativeMask&transform.RelativeRotation != 0 {
					rotTransaction.TriggerChangeListener(child)
				}
			}
		}
		t.logger.Warn(ecs.FlushMany(posTransaction, rotTransaction))
	}

	onSizeUpsert := func(ei []ecs.EntityID) {
		posTransaction := t.posArray.Transaction()
		sizeTransaction := t.posArray.Transaction()
		for _, entity := range ei {
			parentObj := t.hierarchyTransaction.GetObject(entity)
			for _, child := range parentObj.Children().GetIndices() {
				mask, err := t.parentMaskArray.GetComponent(child)
				if err != nil {
					continue
				}
				if mask.RelativeMask&transform.RelativePos != 0 {
					posTransaction.TriggerChangeListener(child)
				}
				if mask.RelativeMask&transform.RelativeSizeXYZ != 0 {
					sizeTransaction.TriggerChangeListener(child)
				}
			}
		}
		t.logger.Warn(ecs.FlushMany(posTransaction, sizeTransaction))
	}

	t.posArray.OnChange(onPosUpsert)
	t.pivotPointArray.OnChange(onPosUpsert)
	t.parentPivotPointArray.OnChange(onPosUpsert)
	t.rotationArray.OnChange(onRotUpsert)
	t.sizeArray.OnChange(onSizeUpsert)

	// t.posArray.OnAdd(onPosUpsert)
	// t.pivotPointArray.OnAdd(onPosUpsert)
	// t.parentPivotPointArray.OnAdd(onPosUpsert)
	// t.rotationArray.OnAdd(onRotUpsert)
	// t.sizeArray.OnAdd(onSizeUpsert)
}
