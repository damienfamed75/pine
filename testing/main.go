package main

import (
	"image"
	"image/color"
	"image/draw"
	"path/filepath"

	"github.com/damienfamed75/pine/view"

	"github.com/fogleman/fauxgl"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/oakmound/oak"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
)

func main() {
	oak.SetupConfig.Screen = oak.Screen{
		Width:  1600,
		Height: 900,
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

	// camera *birch.Camera
	camera *view.Camera

	r render.Renderable
}

// NewHello initializes the default values of this scene.
func NewHello() *HelloScene {
	aspect := float64(oak.ScreenWidth) / float64(oak.ScreenHeight)
	return &HelloScene{
		modelPath:   filepath.Join("model", "dwarf.obj"),
		texturePath: "dwarf_diffuse.png",
		// camera: birch.NewStaticCamera(birch.NewVertex(1, 0.75, 1), birch.Vertex{}, birch.Vertex{}, 100),
		camera: view.NewCamera(mgl64.Vec3{1, 0.75, 1}, 70, aspect),
	}
}

// Start is the initializer stage right when the scene is loaded into oak.
func (h *HelloScene) Start(string, interface{}) {
	// For rendering I'm using an N64 obj loader...
	// If this was a practical demo, then I'd make a new obj loader.
	// r, err := birch.NewRender(
	// 	h.camera,
	// 	h.modelPath,
	// 	h.texturePath,
	// 	oak.ScreenWidth,
	// 	oak.ScreenHeight,
	// )

	// if err != nil {
	// 	// Use the oak logger to exit and log the error.
	// 	dlog.Error(err)
	// 	return
	// }

	// Set the renderable object.
	// h.r = r
	obj, err := NewOBJ(h.camera, "dwarf.obj", oak.ScreenWidth, oak.ScreenHeight)
	if err != nil {
		dlog.Error(err)
	}
	h.r = obj

	// Render the 3D model.
	render.Draw(obj)
}

// Loop returns whether this scene should continue or end.
// By always returning true, it indicates that the scene should never stop looping.
func (h *HelloScene) Loop() bool {
	// h.camera.Update()
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

type Obj struct {
	*render.Sprite
	image.Image

	mesh *fauxgl.Mesh
	ctx  *fauxgl.Context
	cam  *view.Camera
}

func NewOBJ(cam *view.Camera, path string, w, h int) (*Obj, error) {
	obj := &Obj{}

	mesh, err := fauxgl.LoadOBJ(path)
	if err != nil {
		return nil, err
	}

	mesh.BiUnitCube()
	mesh.SmoothNormalsThreshold(fauxgl.Radians(30))

	context := fauxgl.NewContext(w, h)

	obj.ctx = context
	obj.mesh = mesh
	obj.cam = cam
	obj.Sprite = render.NewEmptySprite(0, 0, w, h)

	return obj, nil
}

func (o *Obj) Draw(buff draw.Image) {
	o.DrawOffset(buff, 0, 0)
}

func (o *Obj) DrawOffset(buff draw.Image, xOff, yOff float64) {
	mat := o.cam.GetViewProjection()
	fmat := fauxgl.Matrix{
		X00: mat[0],
		X01: mat[1],
		X02: mat[2],
		X03: mat[3],
		X10: mat[4],
		X11: mat[5],
		X12: mat[6],
		X13: mat[7],
		X20: mat[8],
		X21: mat[9],
		X22: mat[10],
		X23: mat[11],
		X30: mat[12],
		X31: mat[13],
		X32: mat[14],
		X33: mat[15],
	}
	o.ctx.Shader = fauxgl.NewSolidColorShader(fmat, fauxgl.Color{R: 255, A: 255})
	o.Image = o.ctx.Image()
	o.ctx.DrawMesh(o.mesh)
	buff = o
	o.Sprite.DrawOffset(buff, xOff, yOff)
}

func (o *Obj) Set(x, y int, c color.Color) {

}
