package view

import (
	"image"
	"image/draw"
	"math"

	"github.com/damienfamed75/pine/tdraw"
	"github.com/go-gl/mathgl/mgl64"
)

// Draw calls DrawOffset at 0 offset.
func (m *Model) Draw(buff draw.Image) {
	m.DrawOffset(buff, 0, 0)
}

// DrawOffset will simply render the model and then offset it on the window
// by the provided x and y offsets.
func (m *Model) DrawOffset(buff draw.Image, xOff, yOff float64) {
	// Get the boundaries of the model's sprite.
	// This should be the width and height assigned.
	bounds := m.Sprite.GetRGBA().Bounds()
	// Reset the sprite's RGBA values.
	m.Sprite.SetRGBA(image.NewRGBA(bounds))

	// Get the sprite's width and height.
	spriteWidth := bounds.Max.X
	spriteHeight := bounds.Max.Y

	// Setup a zbuffer so we know what pixels we should draw and which ones
	// are behind others we have already drawn. Initialize all values in the
	// buffer to be as far back as possible (-math.MaxFloat64)
	zbuff := make([][]float64, spriteWidth)
	for i := range zbuff {
		zbuff[i] = make([]float64, spriteHeight)
		for j := range zbuff[i] {
			zbuff[i][j] = -math.MaxFloat64
		}
	}

	// Rotation gets applied to the camera items to emulate that the camera
	// is viewing the object from a different angle.
	eye := m.quat.Rotate(m.camera.position)    // rotate position.
	up := m.quat.Rotate(m.camera.up)           // rotate upward rotation.
	forward := m.quat.Rotate(m.camera.forward) // rotate forward rotation.

	// Get the camera's projection matrix.
	// This allows us to create perspective when drawing the object.
	proj := m.camera.GetViewProjection()

	// If I'm able to remove these values that'd be create.
	//
	// These are used to calculate some simple projections of the model's
	// triangles.
	z := eye.Sub(forward).Normalize()
	x := up.Cross(z).Normalize()
	y := z.Cross(x)

	// For every triangle in the model.
	//
	// Loop loops every three vertices, because we need three vertices to build
	// a triangle to render.
	for i := 0; i < len(m.outVertices); i += 3 {
		var (
			mvert, mnrm, mtex tdraw.Triangle
		)

		// Vertex Normals.
		mnrm = tdraw.Triangle{
			A: m.outNormals[i],
			B: m.outNormals[i+1],
			C: m.outNormals[i+2],
		}.ViewNrm(x, y, z)

		// Model Coordinates.
		mvert = tdraw.Triangle{
			A: m.outVertices[i],
			B: m.outVertices[i+1],
			C: m.outVertices[i+2],
		}.ViewTri(x, y, z, eye)

		// Texture Coordinates.
		mtex = tdraw.Triangle{
			A: m.outUVs[i],
			B: m.outUVs[i+1],
			C: m.outUVs[i+2],
		}

		// Perspective Vertices.
		vew := tdraw.Triangle{
			A: mgl64.Project(
				mvert.A, m.transform, proj, 0, 0, spriteWidth, spriteHeight),
			B: mgl64.Project(
				mvert.B, m.transform, proj, 0, 0, spriteWidth, spriteHeight),
			C: mgl64.Project(
				mvert.C, m.transform, proj, 0, 0, spriteWidth, spriteHeight),
		}

		// Draw the triangles into the buffer.
		tdraw.TDraw(
			m.Sprite.GetRGBA(),
			zbuff,
			vew,
			mnrm,
			mtex,
			m.textureData,
		)
	}

	// Let oak render the buffer onto the window.
	m.Sprite.DrawOffset(buff, xOff, yOff)
}

func getFaceVertices(mod *Model, indices []uint, tmpVerts, out []mgl64.Vec3) {
	for i := range indices {
		idx := indices[i]

		out = append(out, tmpVerts[idx-1])
	}
}

func mulQuatTriangle(t tdraw.Triangle, quat mgl64.Quat) tdraw.Triangle {
	return tdraw.Triangle{
		A: quat.Rotate(t.A),
		B: quat.Rotate(t.B),
		C: quat.Rotate(t.C),
	}
}
