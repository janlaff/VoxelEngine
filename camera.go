package main

import (
	"github.com/Ferguzz/glam/math"
	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	programId uint32
	position  mgl32.Vec3
	direction mgl32.Vec3
}

func (c *Camera) SetModel(model mgl32.Mat4) {
	modelLocation := gl.GetUniformLocation(c.programId, gl.Str("model\x00"))
	gl.Uniform4fv(modelLocation, 1, &model[0])
}

func (c *Camera) SetPosition(position mgl32.Vec3) {
	c.position = position

	cameraPositionLocation := gl.GetUniformLocation(c.programId, gl.Str("cameraPosition\x00"))
	gl.Uniform3fv(cameraPositionLocation, 1, &c.position[0])
}

func (c *Camera) SetOrientation(yaw float32, pitch float32) {
	c.direction = mgl32.Vec3{
		math.Cos(mgl32.DegToRad(yaw)) * math.Cos(mgl32.DegToRad(pitch)),
		math.Sin(mgl32.DegToRad(pitch)),
		math.Sin(mgl32.DegToRad(yaw)) * math.Cos(mgl32.DegToRad(pitch)),
	}.Normalize()

	c.UpdateView()
}

func (c *Camera) LookAt(point mgl32.Vec3) {
	c.direction = point.Sub(c.position)
	c.UpdateView()
}

func (c *Camera) SetScreenSize(screenWidth float32, screenHeight float32) {
	projection := mgl32.Perspective(
		mgl32.DegToRad(45),
		screenWidth/screenHeight,
		0.1,
		100,
	)

	gl.Viewport(0, 0, int32(screenWidth), int32(screenHeight))

	projectionLocation := gl.GetUniformLocation(c.programId, gl.Str("projection\x00"))
	screenWidthLocation := gl.GetUniformLocation(c.programId, gl.Str("screenWidth\x00"))
	screenHeightLocation := gl.GetUniformLocation(c.programId, gl.Str("screenHeight\x00"))

	gl.UniformMatrix4fv(projectionLocation, 1, false, &projection[0])
	gl.Uniform1f(screenWidthLocation, screenWidth)
	gl.Uniform1f(screenHeightLocation, screenHeight)
}

func (c *Camera) UpdateView() {
	target := c.position.Add(c.direction)

	view := mgl32.LookAt(
		c.position.X(),
		c.position.Y(),
		c.position.Z(),
		target.X(),
		target.Y(),
		target.Z(),
		0,
		1,
		0,
	)

	viewLocation := gl.GetUniformLocation(c.programId, gl.Str("view\x00"))
	gl.UniformMatrix4fv(viewLocation, 1, false, &view[0])
}
