package view

import (
	"github.com/go-gl/mathgl/mgl64"
)

// Camera is a 3-Dimensional camera object with a perspective transform to apply
// to the around around us.
//
// TODO more documentation on camera structure.
type Camera struct {
	// if the window is resized then the aspect ratio won't work
	// and then the perspective matrix will break.
	perspective mgl64.Mat4 // The perspective matrix to apply to other matrices.
	transform   mgl64.Mat4
	position    mgl64.Vec3 // The position of the camera in the world.

	// Rotation based vectors.
	forward mgl64.Vec3 // What the camera sees as forward. (typically Z or Y)
	up      mgl64.Vec3 // What the camera sees as up. (typically Y or Z)
}

// NewExplicitCamera returns a new camera with a transform perspective matrix.
// The parameters include:
// 1. pos - position of the camera in the world space
// 2. forward - what the camera perceives as forward
// 3. up - what the camera perceives as up
// 4. fovy - the camera's field of vision
// 5. aspect - the aspect ratio of the viewport. (width ÷ height)
// 6. zNear - the nearest z position the camera can see
// 7. zFar - the furthest z position the camera can see
//
// Warning: To not set zNear and zFar too far apart or else there could be some
// floating point precision errors that arise.
func NewExplicitCamera(pos, forward, up mgl64.Vec3, fovy, aspect, zNear, zFar float64) *Camera {
	return &Camera{
		perspective: mgl64.Perspective(fovy, aspect, zNear, zFar),
		position:    pos,
		forward:     forward,
		up:          up,
		transform: mgl64.LookAtV(
			pos, pos.Add(forward), up,
		),
	}
}

// NewCamera is a somewhat defaulted value camera.
func NewCamera(pos mgl64.Vec3, fovy, aspect float64) *Camera {
	return NewExplicitCamera(
		pos,                 // Position of the camera.
		mgl64.Vec3{0, 0, 1}, // Z axis is what we perceive is forward.
		mgl64.Vec3{0, 1, 0}, // Y is what we perceive is up.
		fovy,                // Field of vision.
		aspect,              // Aspect ratio.
		0.01,                // The closest we can see.
		1000,                // The furthest we can see.
	)
}

func (c *Camera) Rotate(angle float64) {
	rz := mgl64.Rotate3DZ(angle)
	c.forward = rz.Mul3x1(c.forward)
}

func (c *Camera) GetForwardRot() mgl64.Vec3 {
	return c.forward
}

func (c *Camera) GetUpRot() mgl64.Vec3 {
	return c.up
}

func (c *Camera) GetTransform() mgl64.Mat4 {
	return c.transform
}

// GetViewProjection gets a transform matrix of the perspective matrix
// to apply to the objects around us.
func (c *Camera) GetViewProjection() mgl64.Mat4 {
	// LookAtV generates a transform matrix from world space into eye space.
	// The eye being the camera.
	//
	// Params of LookAtV:
	// 1. where we are
	// 2. What we are looking at. That being what we perceive is forward.
	// 3. What we perceive is up.
	//
	// We multiply the eye space by our perspective transform matrix to apply
	// the perspective to the world around us.
	return c.perspective.Mul4(
		mgl64.LookAtV(
			c.position,
			c.position.Add(c.forward),
			c.up,
		),
	)
}

func (c *Camera) GetPerspective() mgl64.Mat4 {
	return c.perspective
}

func (c *Camera) GetPosition() mgl64.Vec3 {
	return c.position
}
