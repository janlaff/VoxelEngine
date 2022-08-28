package main

import (
	_ "embed"
	"fmt"
	"github.com/Ferguzz/glam/math"
	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"runtime"
)

//go:embed assets/shaders/voxels.vert
var vertexShaderSource string

//go:embed assets/shaders/voxels.frag
var fragmentShaderSource string

type DeltaTimer struct {
	lastFrame float64
	deltaTime float64
}

func NewDeltaTimer() DeltaTimer {
	var dt DeltaTimer
	dt.lastFrame = glfw.GetTime()
	dt.deltaTime = 0
	return dt
}

func (dt *DeltaTimer) Update() {
	currentFrame := glfw.GetTime()
	dt.deltaTime = currentFrame - dt.lastFrame
	dt.lastFrame = currentFrame
}

type FpsCounter struct {
	lastFrame  float64
	frameCount uint32
}

func NewFpsCounter() FpsCounter {
	var f FpsCounter
	f.frameCount = 0
	f.lastFrame = glfw.GetTime()
	return f
}

func (fps *FpsCounter) Update() {
	currentFrame := glfw.GetTime()
	fps.frameCount++

	if currentFrame-fps.lastFrame >= 1 {
		fmt.Println("Fps: ", fps.frameCount)
		fps.frameCount = 0
		fps.lastFrame = currentFrame
	}
}

func main() {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	window := openWindow()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	vertexShader, err := compileShader(vertexShaderSource+"\x00", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource+"\x00", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program, err := linkProgram([]uint32{vertexShader, fragmentShader})
	if err != nil {
		panic(err)
	}

	gl.ClearColor(0.67, 0.84, 0.9, 1.0)
	gl.UseProgram(program)

	var chunk Chunk
	chunk.data[0] = true
	chunk.data[1] = true
	meshData := chunk.CreateMesh()
	//meshData.AddCube(mgl32.Vec3{0, 0, 0})

	var vertexBuffer uint32
	gl.GenBuffers(1, &vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(meshData.vertices), gl.Ptr(meshData.vertices), gl.STATIC_DRAW)

	var elementBuffer uint32
	gl.GenBuffers(1, &elementBuffer)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, elementBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(meshData.triangles), gl.Ptr(meshData.triangles), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, elementBuffer)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.FRONT)

	var camera Camera
	camera.programId = program
	camera.SetModel(mgl32.Ident4())
	camera.LookAt(mgl32.Vec3{0, 0, 0})
	camera.SetScreenSize(800, 600)

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		camera.SetScreenSize(float32(width), float32(height))
	})

	dt := NewDeltaTimer()
	fps := NewFpsCounter()

	for !window.ShouldClose() {
		dt.Update()
		fps.Update()

		rotationSpeed := glfw.GetTime() * 0.3
		radius := float32(4)
		camX := math.Sin(float32(rotationSpeed)) * radius
		camZ := math.Cos(float32(rotationSpeed)) * radius
		camera.SetPosition(mgl32.Vec3{
			camX, 1, camZ,
		})
		camera.LookAt(mgl32.Vec3{
			0, 0, 0,
		})

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)
		gl.BindVertexArray(vao)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, elementBuffer)
		gl.DrawElements(gl.TRIANGLES, int32(len(meshData.triangles)), gl.UNSIGNED_INT, nil)
		glfw.PollEvents()
		window.SwapBuffers()

		if window.GetKey(glfw.KeyQ) == glfw.Press {
			break
		}
	}
}

func openWindow() *glfw.Window {
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 4)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	window, err := glfw.CreateWindow(800, 600, "VoxelEngine", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	return window
}
