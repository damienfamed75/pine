package view

import (
	"image"

	"github.com/oakmound/oak/render"
)

type Mesh struct {
	// a render.Sprite has a position and a buffer of image data which
	// it uses to draw to the screen at that position.
	*render.Sprite
	// the textureData is the local texture file (.bmp in the original, .png in this version)
	// that is referred to to color each triangle face
	textureData *image.RGBA
}

func LoadObj(objFile, texFile string, w, h int) error {
	// fobj, err := os.Open(objfile)
	// if err != nil {
	// 	return err
	// }
	// tex, err := render.LoadSprite("model", texfile)
	// if err != nil {
	// 	return err
	// }

	// mgl64.

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

	return nil
}
