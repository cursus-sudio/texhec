package client

// import "engine/services/ecs"
//
// // on event from config start prediction recording
// func NewStartPredictionSystem(
// 	toolFactory ecs.ToolFactory[Tool],
// ) ecs.SystemRegister {
// 	return ecs.NewSystemRegister(func(w ecs.World) error {
// 		tool := toolFactory.Build(w)
// 		for _, listen := range tool.ListenToEvents {
// 			listen(w.EventsBuilder(), tool.BeforeInternalEvent)
// 		}
// 		return nil
// 	})
// }
//
// //
//
// // on event from config stop prediction recording
// func NewStopPredictionSystem(toolFactory ecs.ToolFactory[Tool]) ecs.SystemRegister {
// 	return ecs.NewSystemRegister(func(w ecs.World) error {
// 		tool := toolFactory.Build(w)
// 		for _, listen := range tool.ListenToEvents {
// 			listen(w.EventsBuilder(), tool.AfterInternalEvent)
// 		}
// 		return nil
// 	})
// }
