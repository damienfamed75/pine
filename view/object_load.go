package view

import (
	"bufio"
	"fmt"
	"os"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/oakmound/oak/render"
)

// LoadObj loads a .obj file into memory, loading all its information
// including the texture for the .obj file.
//
// v - vertices
// vn - vertex normalized
// vt - vertex texture coordinate
// f - faces (triangles)
// mtl files are for another time for now.
// f 1/13/4 51/13/5 2/42/26
//				  3rd coord
//        2nd coord
// 1st coord
func LoadObj(objFile, texFile string, w, h int, camera *Camera) (*Model, error) {
	fobj, err := os.Open(objFile)
	if err != nil {
		return nil, err
	}
	defer fobj.Close()

	tex, err := render.LoadSprite("model", texFile)
	if err != nil {
		return nil, err
	}

	mod := &Model{
		// Raw texture data from pixel to pixel.
		textureData: tex.GetRGBA(),
		// Empty sprite that has an assigned width and height.
		Sprite: render.NewEmptySprite(0, 0, w, h),
		// Enough data to render the object.
		camera: camera,
		// quat: mgl64.QuatRotate(mgl64.DegToRad(0), mgl64.Vec3{0, 1, 0}),
		transform: mgl64.Scale3D(1, 1, 1).Mul4(mgl64.Translate3D(0, 0, 0)),
	}

	// Get the raw texture data from pixel to pixel.
	mod.textureData = tex.GetRGBA()
	mod.Sprite = render.NewEmptySprite(0, 0, w, h)
	mod.camera = camera
	// quat := mgl64.QuatIdent().Rotate(mgl64.Vec3{
	// 	0, mgl64.DegToRad(45), 0,
	// })
	mod.quat = mgl64.QuatRotate(mgl64.DegToRad(0), mgl64.Vec3{0, 1, 0})
	// mod.transform = mgl64.Translate3D(0, 0, 0).
	// 	Mul4(quat.Mat4()).
	// 	Mul4(mgl64.Scale3D(1, 1, 1))
	// mod.transform = mgl64.Translate3D(0, 0, -2).
	// 	Mul4(mgl64.HomogRotate3D(mgl64.DegToRad(45), mgl64.Vec3{0, 1, 0})).
	// 	Mul4(mgl64.Scale3D(1, 1, 1))
	mod.transform = mgl64.Scale3D(1, 1, 1).
		Mul4(mgl64.Translate3D(0, 0, 0))

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

	return mod, nil
}
