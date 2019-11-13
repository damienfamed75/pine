package view

import (
	"image"
	"sync"

	"github.com/damienfamed75/pine/tdraw"
	"github.com/go-gl/mathgl/mgl64"
)

// workerPackage contains all the data necessary to render pixels on to the
// given buffer.
type workerPackage struct {
	x, y, z mgl64.Vec3
	eye     mgl64.Vec3

	spriteDimensions *image.Point

	transform mgl64.Mat4
	proj      mgl64.Mat4

	outUVs      []mgl64.Vec3
	outVertices []mgl64.Vec3
	outNormals  []mgl64.Vec3

	zbuff [][]float64

	spriteRGBA  *image.RGBA
	textureData *image.RGBA
}

// For every triangle in the model.
//
// Loop loops every three vertices, because we need three vertices to build
// a triangle to render.
func drawingWorker(pkg *workerPackage, indices chan int, wg *sync.WaitGroup) {
	var (
		mvert, mnrm, mtex tdraw.Triangle
	)

	for i := range indices {

		// Vertex Normals.
		mnrm = tdraw.Triangle{
			A: pkg.outNormals[i],
			B: pkg.outNormals[i+1],
			C: pkg.outNormals[i+2],
		}.ViewNrm(pkg.x, pkg.y, pkg.z)

		// Model Coordinates.
		mvert = tdraw.Triangle{
			A: pkg.outVertices[i],
			B: pkg.outVertices[i+1],
			C: pkg.outVertices[i+2],
		}.ViewTri(pkg.x, pkg.y, pkg.z, pkg.eye)

		// Texture Coordinates.
		mtex = tdraw.Triangle{
			A: pkg.outUVs[i],
			B: pkg.outUVs[i+1],
			C: pkg.outUVs[i+2],
		}

		// Perspective Vertices.
		vew := tdraw.Triangle{
			A: mgl64.Project(
				mvert.A, pkg.transform, pkg.proj, 0, 0, pkg.spriteDimensions.X, pkg.spriteDimensions.Y),
			B: mgl64.Project(
				mvert.B, pkg.transform, pkg.proj, 0, 0, pkg.spriteDimensions.X, pkg.spriteDimensions.Y),
			C: mgl64.Project(
				mvert.C, pkg.transform, pkg.proj, 0, 0, pkg.spriteDimensions.X, pkg.spriteDimensions.Y),
		}

		// Draw the triangles into the buffer.
		tdraw.TDraw(
			pkg.spriteRGBA,
			pkg.zbuff,
			vew,
			mnrm,
			mtex,
			pkg.textureData,
		)
	}

	wg.Done()
}
