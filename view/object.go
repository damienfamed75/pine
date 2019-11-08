package view

import (
	"bufio"
	"fmt"
	"image"
	"image/draw"
	"math"
	"os"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/oakmound/oak/render"
)

type Model struct {
	// a render.Sprite has a position and a buffer of image data which
	// it uses to draw to the screen at that position.
	*render.Sprite
	// the textureData is the local texture file (.bmp in the original, .png in this version)
	// that is referred to to color each triangle face
	textureData *image.RGBA

	outVertices []mgl64.Vec3
	outUVs      []mgl64.Vec3
	outNormals  []mgl64.Vec3

	transform mgl64.Mat4

	camera *Camera
}

func LoadObj(objFile, texFile string, w, h int, camera *Camera) (*Model, error) {
	mod := &Model{}

	fobj, err := os.Open(objFile)
	if err != nil {
		return nil, err
	}
	defer fobj.Close()

	tex, err := render.LoadSprite("model", texFile)
	if err != nil {
		return nil, err
	}

	// Get the raw texture data from pixel to pixel.
	mod.textureData = tex.GetRGBA()
	mod.Sprite = render.NewEmptySprite(0, 0, w, h)

	var (
		uvIndices     []uint
		vertexIndices []uint
		normalIndices []uint

		tmpUVs      []mgl64.Vec3
		tmpVertices []mgl64.Vec3
		tmpNormals  []mgl64.Vec3
	)

	scanner := bufio.NewScanner(fobj)

	for scanner.Scan() {
		var (
			vertex struct{ x, y, z float64 }
		)

		line := scanner.Text()

		if len(line) < 2 {
			continue
		}
		if line[0] == 'v' && line[1] == 'n' {
			// vertex normals.
			fmt.Sscanf(line, "vn %f %f %f", &vertex.x, &vertex.y, &vertex.z)
			tmpNormals = append(tmpNormals, mgl64.Vec3{
				vertex.x, vertex.y, vertex.z,
			})
		} else if line[0] == 'v' && line[1] == 't' {
			// vertex texture coordinates.
			// Most of the time an obj file will not have a Z point, but we'll
			// include it anyway in the case that an obj file actually uses it.
			fmt.Sscanf(line, "vt %f %f %f", &vertex.x, &vertex.y, &vertex.z)
			tmpUVs = append(tmpUVs, mgl64.Vec3{
				vertex.x, vertex.y, vertex.z,
			})
		} else if line[0] == 'v' {
			// vertices
			fmt.Sscanf(line, "v %f %f %f", &vertex.x, &vertex.y, &vertex.z)
			tmpVertices = append(tmpVertices, mgl64.Vec3{
				vertex.x, vertex.y, vertex.z,
			})
		} else if line[0] == 'f' {
			var (
				uvIndex     [3]uint
				vertexIndex [3]uint
				normalIndex [3]uint
			)

			fmt.Sscanf(line, "f %d/%d/%d %d/%d/%d %d/%d/%d",
				&vertexIndex[0], &uvIndex[0], &normalIndex[0],
				&vertexIndex[1], &uvIndex[1], &normalIndex[1],
				&vertexIndex[2], &uvIndex[2], &normalIndex[2],
			)

			uvIndices = append(uvIndices,
				uvIndex[0], uvIndex[1], uvIndex[2])
			vertexIndices = append(vertexIndices,
				vertexIndex[0], vertexIndex[1], vertexIndex[2])
			normalIndices = append(normalIndices,
				normalIndex[0], normalIndex[1], normalIndex[2])
		}
	}

	// Looping through the faces and getting their according vertices.
	for i := range vertexIndices {
		vertIdx := vertexIndices[i]
		// The -1 is because OBJ files for arrays start at 1 not 0.
		// So to compensate for Golang we are subtracting the index by one.
		mod.outVertices = append(mod.outVertices, tmpVertices[vertIdx-1])
	}
	for i := range uvIndices {
		vertIdx := uvIndices[i]
		// The -1 is because OBJ files for arrays start at 1 not 0.
		// So to compensate for Golang we are subtracting the index by one.
		mod.outUVs = append(mod.outUVs, tmpUVs[vertIdx-1])
	}
	for i := range normalIndices {
		vertIdx := normalIndices[i]
		// The -1 is because OBJ files for arrays start at 1 not 0.
		// So to compensate for Golang we are subtracting the index by one.
		mod.outNormals = append(mod.outNormals, tmpNormals[vertIdx-1])
	}

	// figure out loading an obj file w/ opengl matrices
	// v - vertices
	// vn - vertex normalized
	// vt - vertex texture coordinate
	// f - faces (triangles)
	// mtl files are for another time for now.
	// f 1/13/4 51/13/5 2/42/26
	//				  3rd coord
	//        2nd coord
	// 1st coord
	//
	//

	// figure out rendering with the opengl matrices instead
	// of the last solution.

	return mod, nil
}

func (m *Model) Draw(buff draw.Image) {
	m.DrawOffset(buff, 0, 0)
}

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

	proj := m.camera.GetViewProjection()

	for i := 0; i < len(m.outVertices); i++ {
		m.outNormals[i] = m.outNormals[i].Mul(proj.Det())
		m.outVertices[i] = m.outVertices[i].Mul(proj.Det())
		m.outNormals[i] = m.outNormals[i].Mul(proj.Det())
	}
}

func getFaceVertices(mod *Model, indices []uint, tmpVerts, out []mgl64.Vec3) {
	for i := range indices {
		idx := indices[i]

		out = append(out, tmpVerts[idx-1])
	}
}
