package view

import (
	"image"

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

	// quat represents the model's rotation.
	// the quaternion isn't directly applied to the transform, but instead
	// applied to the camera's viewing position and angles to create
	// the illusion of full 3D rotation.
	quat mgl64.Quat
	// transform represents the model's position and scale.
	transform mgl64.Mat4
	angle     float64

	camera *Camera
}

// Rotate resets the rotation of the model to what is provided.
// If you wish to rotate based on the current rotation then please refer to
// RotateExisting instead.
func (m *Model) Rotate(angle float64, axis mgl64.Vec3) {
	m.angle = angle
	m.quat = mgl64.QuatRotate(angle, axis)
}

func (m *Model) RotateExisting(angle float64, axis mgl64.Vec3) {
	m.angle += angle
	m.quat = mgl64.QuatRotate(m.angle, axis)
}
