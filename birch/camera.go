package birch

// DEPRECATED

// import (
// 	"math"

// 	"github.com/oakmound/oak/mouse"
// )

// type Camera struct {
// 	Position   Vertex
// 	Target     Vertex
// 	Up         Vertex
// 	Fovy       float64
// 	DeltaMouse mouse.Event

// 	// If the camera is movable, then have lastmouse event.
// 	sensitivity float64
// 	static      bool
// }

// func NewMovableCamera(pos, target, up Vertex, fovy, sensitivity float64) *Camera {
// 	if sensitivity == 0.0 {
// 		sensitivity = .005 // default if the sensitivity is zero.
// 	}

// 	return &Camera{
// 		Position:    pos,
// 		Target:      target,
// 		Up:          up,
// 		Fovy:        fovy,
// 		sensitivity: sensitivity,
// 		static:      false,
// 	}
// }

// func NewStaticCamera(pos, target, up Vertex, fovy float64) *Camera {
// 	return &Camera{
// 		Position: pos,
// 		Target:   target,
// 		Up:       up,
// 		Fovy:     fovy,
// 		static:   true,
// 	}
// }

// func (c *Camera) Update() {
// 	if !c.static && mouse.LastEvent != c.DeltaMouse {
// 		mouseXt := mouse.LastEvent.X() * c.sensitivity
// 		mouseYt := mouse.LastEvent.Y() * c.sensitivity

// 		// The math here must be incorrect, because it acts wonky.
// 		c.Position = Vertex{
// 			math.Sin(mouseXt), math.Sin(mouseYt), math.Cos(mouseXt),
// 		}
// 	}
// }
