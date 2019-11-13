package main

import (
	"path/filepath"

	"github.com/damienfamed75/pine/view"
	"github.com/go-gl/mathgl/mgl64"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
)

const (
	screenWidth  = 1600
	screenHeight = 900
)

func main() {
	oak.SetupConfig.Screen = oak.Screen{
		Width:  screenWidth,
		Height: screenHeight,
	}

	oak.SetupConfig.Title = "Hello"
	oak.SetupConfig.BatchLoad = true
	oak.SetupConfig.Assets = oak.Assets{
		AssetPath: "/",
		ImagePath: "model",
	}

	hello := NewHello()

	// idea...
	// new scene
	// newobj(scene, params...)
	// -> scene is used to obtain camera and then obj gets stored in scene

	oak.Add("hello", hello.Start, hello.Loop, hello.End)
	oak.Init("hello")
}

// HelloScene is a testing scene just to render a 3D model of a dwarf...
type HelloScene struct {
	modelPath   string
	texturePath string

	// camera *birch.Camera
	camera *view.Camera

	r *view.Model
	// r render.Renderable
}

// NewHello initializes the default values of this scene.
func NewHello() *HelloScene {
	aspect := float64(screenWidth) / float64(screenHeight)

	return &HelloScene{
		modelPath:   filepath.Join("model", "dwarf.obj"),
		texturePath: "dwarf_diffuse.png",
		camera:      view.NewCamera(mgl64.Vec3{1, 0.75, 1}, mgl64.DegToRad(90), aspect),
	}
}

// Start is the initializer stage right when the scene is loaded into oak.
func (h *HelloScene) Start(string, interface{}) {
	// Load an obj file using the new object loader in view package.
	// Also load in the texture and store the necessary values to render it on
	// the screen.
	r, err := view.LoadObj(
		h.modelPath,
		h.texturePath,
		oak.ScreenWidth,
		oak.ScreenHeight,
		h.camera,
	)
	if err != nil {
		// Use the oak logger to exit and log the error.
		dlog.Error(err)
		return
	}

	// Set the renderable object.
	h.r = r

	// Render the 3D model.
	render.Draw(r)
}

// Loop returns whether this scene should continue or end.
// By always returning true, it indicates that the scene should never stop looping.
func (h *HelloScene) Loop() bool {
	h.r.AddRotation(mgl64.RadToDeg(.00005), mgl64.Vec3{0, 1, 0})
	return true
}

// End is never called, but were it called it would
// end the gel scene and start the gel scene again anew
// the return values are 1) the scene to go next and 2)
// any settings that should be applied when transitioning
// to the next scene, in this case none.
func (h *HelloScene) End() (string, *scene.Result) {
	return "hello", nil
}
