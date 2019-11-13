package view

import (
	"image"
	"image/draw"
	"math"
	"sync"

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

	// Setup a zbuffer so we know what pixels we should draw and which ones
	// are behind others we have already drawn. Initialize all values in the
	// buffer to be as far back as possible (-math.MaxFloat64)
	zbuff := make([][]float64, bounds.Max.X)
	for i := range zbuff {
		zbuff[i] = make([]float64, bounds.Max.Y)
		for j := range zbuff[i] {
			zbuff[i][j] = -math.MaxFloat64
		}
	}

	// Rotation gets applied to the camera items to emulate that the camera
	// is viewing the object from a different angle.
	eye := m.quat.Rotate(m.camera.GetPosition())            // position.
	up := m.quat.Rotate(m.camera.GetUpRotation())           // upward rotation.
	forward := m.quat.Rotate(m.camera.GetForwardRotation()) // forward rotation.

	// Get the camera's projection matrix.
	// This allows us to create perspective when drawing the object.
	proj := m.camera.GetViewProjection()

	transform := m.GetTransform()

	// If I'm able to remove these values that'd be create.
	//
	// These are used to calculate some simple projections of the model's
	// triangles.
	z := eye.Sub(forward).Normalize()
	x := up.Cross(z).Normalize()
	y := z.Cross(x)

	var (
		wg      sync.WaitGroup
		indices = make(chan int)
		pkg     = workerPackage{
			outUVs:      m.outUVs,
			outVertices: m.outVertices,
			outNormals:  m.outNormals,
			zbuff:       zbuff,
			spriteRGBA:  m.Sprite.GetRGBA(),
			textureData: m.textureData,
			x:           x,
			y:           y,
			z:           z,
			eye:         eye,
			// Get the sprite's width and height.
			spriteDimensions: &image.Point{X: bounds.Max.X, Y: bounds.Max.Y},
			proj:             proj,
			transform:        transform,
		}
	)

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go drawingWorker(&pkg, indices, &wg)
	}

	for i := 0; i < len(m.outVertices); i += 3 {
		indices <- i
	}

	// Close the indices channel to signal that the workers may stop.
	close(indices)

	// Wait for all the workers to finish their work.
	wg.Wait()

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
