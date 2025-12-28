package cameratool

import (
	"engine/modules/camera"
	"engine/modules/collider"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/media/window"
	"reflect"

	"github.com/go-gl/mathgl/mgl32"
)

type ToolFactory interface {
	Register(
		reflect.Type,
		func(camera.World, camera.CameraTool) ProjectionData,
	)
	camera.ToolFactory
}

type ProjectionData struct {
	Mat4     func(entity ecs.EntityID) mgl32.Mat4
	ShootRay func(entity ecs.EntityID, mousePos window.MousePos) collider.Ray
}

type toolFactory struct {
	window window.Api

	projectionIDs map[reflect.Type]projectionID
	projections   datastructures.SparseArray[projectionID, func(camera.World, camera.CameraTool) ProjectionData]
}

func NewCameraResolverFactory(window window.Api) ToolFactory {
	return &toolFactory{
		window: window,

		projectionIDs: make(map[reflect.Type]projectionID),
		projections:   datastructures.NewSparseArray[projectionID, func(camera.World, camera.CameraTool) ProjectionData](),
	}
}

func (f *toolFactory) Register(
	componentType reflect.Type,
	data func(camera.World, camera.CameraTool) ProjectionData,
) {
	if _, ok := f.projectionIDs[componentType]; ok {
		return
	}
	i := projectionID(len(f.projections.GetIndices()))
	f.projectionIDs[componentType] = i
	f.projections.Set(i, data)
}

func (f *toolFactory) Build(world camera.World) camera.CameraTool {
	if t, ok := ecs.GetGlobal[tool](world); ok {
		return &t
	}
	t := tool{
		World: world,

		cameraArray:      ecs.GetComponentsArray[camera.Component](world),
		projectionsArray: ecs.GetComponentsArray[projectionComponent](world),

		toolFactory: f,
		projections: datastructures.NewSparseArray[projectionID, ProjectionData](),

		dirtySet: ecs.NewDirtySet(),

		mobileCamera:       ecs.GetComponentsArray[camera.MobileCameraComponent](world),
		cameraLimits:       ecs.GetComponentsArray[camera.CameraLimitsComponent](world),
		viewport:           ecs.GetComponentsArray[camera.ViewportComponent](world),
		normalizedViewport: ecs.GetComponentsArray[camera.NormalizedViewportComponent](world),

		ortho:              ecs.GetComponentsArray[camera.OrthoComponent](world),
		orthoResolution:    ecs.GetComponentsArray[camera.OrthoResolutionComponent](world),
		perspective:        ecs.GetComponentsArray[camera.PerspectiveComponent](world),
		dynamicPerspective: ecs.GetComponentsArray[camera.DynamicPerspectiveComponent](world),
	}

	world.SaveGlobal(t)

	t.projectionsArray.BeforeGet(t.BeforeGet)
	t.cameraArray.AddDirtySet(t.dirtySet)

	for _, id := range f.projections.GetIndices() {
		value, _ := f.projections.Get(id)
		t.projections.Set(id, value(world, t))
	}
	return t
}
