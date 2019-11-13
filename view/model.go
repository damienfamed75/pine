package view

import (
	"image"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/oakmound/oak/render"
)

// Model represents a 3D model in world space.
// This model is a render.Renderable.
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

	scale    mgl64.Mat4
	position mgl64.Mat4
	// transform represents the model's position and scale.
	// transform mgl64.Mat4
	angle float64

	camera *Camera
}

// SetRotation resets the rotation of the model to what is provided.
// If you wish to rotate based on the current rotation then please refer to
// AddRotation instead.
func (m *Model) SetRotation(angle float64, axis mgl64.Vec3) {
	m.angle = angle
	m.quat = mgl64.QuatRotate(angle, axis)
}

// AddRotation adds to the existing rotation axis.
func (m *Model) AddRotation(angle float64, axis mgl64.Vec3) {
	m.angle += angle
	m.quat = mgl64.QuatRotate(m.angle, axis)
}

// GetTransform combines the scale and position to give the transform matrix of
// this model.
func (m *Model) GetTransform() mgl64.Mat4 {
	return m.scale.Mul4(m.position)
}

// GetScale gets the model's scale on its x, y, and z axis.
func (m *Model) GetScale() mgl64.Vec3 {
	return m.scale.Diag().Vec3()
}

// SetScale sets the model's relative scale on the x, y, and z axis
// If you wish to scale based on thge current rotation then please refer to
// AddScale instead.
func (m *Model) SetScale(x, y, z float64) {
	m.scale = mgl64.Scale3D(x, y, z)
}

// AddScale scales the object relative to its current scale.
func (m *Model) AddScale(x, y, z float64) {
	m.scale = m.scale.Add(mgl64.Scale3D(x, y, z))
}

// SetPosition sets the position of the object from 0,0,0.
// If you wish to set the scale based on its current position then please refer
// to AddPosition instead.
func (m *Model) SetPosition(x, y, z float64) {
	m.position = mgl64.Translate3D(x, y, z)
}

// AddPosition sets the model's position relative to its current position.
func (m *Model) AddPosition(x, y, z float64) {
	m.position = m.position.Add(mgl64.Translate3D(x, y, z))
}
