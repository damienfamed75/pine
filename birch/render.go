package birch

import (
	"image"
	"image/draw"
	"math"
	"os"

	"github.com/damienfamed75/pine/view"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/oakmound/oak/mouse"
	"github.com/oakmound/oak/render"
)

// The Render type implements oak/render.Renderable, letting it be drawn
// to an oak screen. Most of the functions of that interface are satisfied
// by embedding *render.Sprite into the struct, but we replace the Draw
// functions, which are called each frame, below.
type Render struct {
	// a render.Sprite has a position and a buffer of image data which
	// it uses to draw to the screen at that position.
	*render.Sprite
	// the textureData is the local texture file (.bmp in the original, .png in this version)
	// that is referred to to color each triangle face
	textureData *image.RGBA

	rotation  mgl64.Mat3
	transform mgl64.Mat4
	position  mgl64.Vec3 // the world space where the object is.

	// tv, tt, and tn need better names
	tv []Triangle // triangle vertices - outvertices
	tt []Triangle // triangle textures - outUVs
	tn []Triangle // triangle normalized - outNormals
	// the last mouse event is stored on the render so that when it changes,
	// the render knows to update what should be drawn.
	lastmouse mouse.Event

	// Local pointer to the game's camera.
	// cam *Camera
	cam *view.Camera
}

// NewRender creates a Render type to be drawn to screen.
// If it fails, it will return nil and the cause for failure.
// The inputs are the object and texture file paths to be loaded
// and the width and height of the render to be drawn.
func NewRender(cam *view.Camera, objfile, texfile string, w, h int) (*Render, error) {
	// func NewRender(cam *Camera, objfile, texfile string, w, h int) (*Render, error) {
	// Open both the object and texture file
	// if either fails to load, return nil and
	// the cause.
	fobj, err := os.Open(objfile)
	if err != nil {
		return nil, err
	}
	tex, err := render.LoadSprite("model", texfile)
	if err != nil {
		return nil, err
	}
	// parse the object file into an Object struct
	obj := oparse(fobj)
	// return a Render built on that object struct and
	// the texture file.
	return &Render{
		Sprite:      render.NewEmptySprite(0, 0, w, h),
		tv:          obj.Tvgen(),
		tt:          obj.Ttgen(),
		tn:          obj.Tngen(),
		textureData: tex.GetRGBA(),
		cam:         cam,
		transform:   mgl64.Ident4(),
		rotation:    mgl64.Rotate3DX(mgl64.DegToRad(90)),
		position:    mgl64.Vec3{0, 0, 0},
		// transform:   mgl64.Ident4().Mul4(mgl64.Translate3D(0, 0.5, 0)),
		// transform:   mgl64.Translate3D(0, 0, 0).Mul4(mgl64.Scale3D(1, 1, 1)),
		// We set lastmouse to have the Reset event so that our
		// equality check in DrawOffset will fail on the first
		// draw frame. This allows for the render to be drawn
		// before any mouse input is recorded.
		lastmouse: mouse.Event{Event: "Reset"},
	}, nil
}

// Draw expects the render to draw itself to the input buffer.
func (r *Render) Draw(buff draw.Image) {
	// To avoid duplicating logic, Draw just calls DrawOffset with 0 offsets.
	r.DrawOffset(buff, 0, 0)
}

func Unit(v mgl64.Vec3) mgl64.Vec3 {
	return v.Mul(1.0 / v.Len())
}

// DrawOffset expects the render to draw itself to the input buffer,
// offset from it's logical coordinates by xOff and yOff for x and y respectively
func (r *Render) DrawOffset(buff draw.Image, xOff, yOff float64) {
	// Reset the backing Sprite's color buffer to be empty, so we avoid
	// smearing with the last drawn frame.
	bounds := r.Sprite.GetRGBA().Bounds()
	r.Sprite.SetRGBA(image.NewRGBA(bounds))
	w := bounds.Max.X
	h := bounds.Max.Y

	// Get the mouse position and scale it down so we can use it for the
	// render's rotation.

	// Set up a zbuffer so we know what pixels we should draw and which ones
	// are behind others we have already drawn. Initialize all values in the
	// buffer to be as far back as possible (-math.MaxFloat64)
	zbuff := make([][]float64, w)
	for i := range zbuff {
		zbuff[i] = make([]float64, h)
		for j := range zbuff[i] {
			zbuff[i][j] = -math.MaxFloat64
		}
	}
	// What the camera perceives is forward.
	forward := r.cam.GetForwardRot()
	// Which way up is, or pointing in the y direction
	up := r.cam.GetUpRot()
	// Where we're looking from
	// x affects distance
	// y is height of the camera
	// z i'm unsure...
	eye := r.cam.GetPosition()
	proj := r.cam.GetViewProjection()

	transform := r.transform
	// transform := r.rotation.Mat4().Mul4(r.transform)
	// position := transform.Mul4x1(r.position.Vec4(0)).Vec3()

	// z is the depth of the model. The higher the value, the more depth.
	// 0 being the lowest depth. Making him very flat & orthographic
	// 1 being the highest depth.
	z := eye.Sub(forward).Normalize()
	x := up.Cross(z).Normalize()
	// x := z.Cross(up.Normalize()).Normalize()
	y := z.Cross(x)

	// For each triangle, draw it
	for i := 0; i < len(r.tv); i++ {
		// Obtain the normal and triangle values
		// from our view
		// (More documentation needed here)
		nrm := r.tn[i].ViewNrm(x, y, z)
		tri := r.tv[i].ViewTri(x, y, z, eye)
		tex := r.tt[i]

		vew := Triangle{
			A: mgl64.Project(tri.A, transform, proj, 0, 0, w, h),
			B: mgl64.Project(tri.B, transform, proj, 0, 0, w, h),
			C: mgl64.Project(tri.C, transform, proj, 0, 0, w, h),
		}

		// fmt.Printf("vew: [%v, %v, %v]\n", vew.A.X(), vew.A.Y(), vew.A.Z())
		// vew := per.Viewport(floatgeom.Point2{float64(w), float64(h)})
		// Actually draw the triangle given the values we've calculated
		TDraw(r.Sprite.GetRGBA(), zbuff, vew, nrm, tex, r.textureData)
	}
	// Instead of handling the drawing ourselves, let the embedded Sprite which
	// we've populated the color buffer of draw itself.
	r.Sprite.DrawOffset(buff, xOff, yOff)
}

func ApplyProj(obj mgl64.Vec3, modelview, proj mgl64.Mat4) mgl64.Vec3 {
	obj4 := obj.Vec4(1)

	vpp := proj.Mul4(modelview).Mul4x1(obj4)
	vpp = vpp.Mul(1 / vpp.W())

	return vpp.Vec3()
}
