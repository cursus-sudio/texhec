package render

// type CullSystem struct {
// 	world ecs.World
// }
//
// func (system CullSystem) IsVisible(entity ecs.EntityId) (bool, error) {
// 	cameras := system.world.GetEntitiesWithComponents(ecs.GetComponentPointerType((*Projection)(nil)))
// 	if len(cameras) != 1 {
// 		return false, projection.ErrWorldShouldHaveOneProjection
// 	}
// 	camera := cameras[0]
// 	if err := system.world.GetComponents(camera, &proj); err != nil {
// 		return false, err
// 	}
// 	if err := system.world.GetComponents(camera, &cameraTransform); err != nil {
// 		return false, err
// 	}
// 	return false, nil
// }
