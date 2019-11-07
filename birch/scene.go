package birch

// gelStart initializes the elements of our scene.
// The arguments to this function are 1) the previous
// scene (if any), and 2) any data passed to this scene
// by the previous scene.
//
// As we don't use the arguments, we don't give the arguments
// variable names.
// func gelStart(string, interface{}) {
// 	// Initialize the Render object which will actually
// 	// do the drawing work in this scene.
// 	r, err := NewRender(
// 		filepath.Join("model", "salesman.obj"),
// 		"salesman.png",
// 		oak.ScreenWidth,
// 		oak.ScreenHeight,
// 	)
// 	// If we failed to create a render object, report that we did,
// 	// and don't try to draw the render.
// 	if err != nil {
// 		dlog.Error(err)
// 		return
// 	}
// 	// We call the render package (not to be confused with the gel.Render type)
// 	// and tell it to draw the local gel.Render at each draw frame.
// 	render.Draw(r)
// }

// // gelLoop returns whether this scene should continue or end.
// // By always returning true, it indicates that the scene should never stop looping.
// func gelLoop() bool {
// 	return true
// }

// // gelEnd is never called, but were it called it would
// // end the gel scene and start the gel scene again anew
// // the return values are 1) the scene to go next and 2)
// // any settings that should be applied when transitioning
// // to the next scene, in this case none.
// func gelEnd() (string, *scene.Result) {
// 	return "gel", nil
// }

// // We set Scene to initialize with gelStart,
// // loop repeatedly will gelLoop, and terminate
// // with gelEnd.
// var Scene = scene.Scene{
// 	gelStart,
// 	gelLoop,
// 	gelEnd,
// }
