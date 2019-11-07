package main

import (
	"path/filepath"

	"github.com/damienfamed75/pine/birch"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
)

func main() {
	oak.SetupConfig.Screen = oak.Screen{
		Width:  1600,
		Height: 900,
		// Width:  800,
		// Height: 450,
	}

	oak.SetupConfig.Assets = oak.Assets{
		AssetPath: "\\",
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

	camera *birch.Camera

	r render.Renderable
}

// NewHello initializes the default values of this scene.
func NewHello() *HelloScene {
	return &HelloScene{
		modelPath:   filepath.Join("model", "dwarf.obj"),
		texturePath: "dwarf_diffuse.png",
		// camera:      birch.NewMovableCamera(birch.NewVertex(1, 1, 1), birch.Vertex{}, birch.Vertex{}, 100, .005),
		camera: birch.NewStaticCamera(birch.NewVertex(1, 0.75, 1), birch.Vertex{}, birch.Vertex{}, 100),
	}
}

// Start is the initializer stage right when the scene is loaded into oak.
func (h *HelloScene) Start(string, interface{}) {
	// For rendering I'm using an N64 obj loader...
	// If this was a practical demo, then I'd make a new obj loader.
	r, err := birch.NewRender(
		h.camera,
		h.modelPath,
		h.texturePath,
		oak.ScreenWidth,
		oak.ScreenHeight,
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
	h.camera.Update()
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
