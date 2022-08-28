package main

import (
	_ "embed"
	"fmt"
	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"runtime"
)

//go:embed assets/shaders/voxels.vert
var vertexShaderSource string

//go:embed assets/shaders/voxels.frag
var fragmentShaderSource string

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

	planeVertices := []float32{
		// front
		-0.5, 0.5, 0.5,
		0.5, 0.5, 0.5,
		0.5, -0.5, 0.5,
		0.5, -0.5, 0.5,
		-0.5, -0.5, 0.5,
		-0.5, 0.5, 0.5,
		// right side
		0.5, 0.5, 0.5,
		0.5, 0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, -0.5, 0.5,
		0.5, 0.5, 0.5,
	}

	var vertexBuffer uint32
	gl.GenBuffers(1, &vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(planeVertices), gl.Ptr(planeVertices), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.UseProgram(program)

	var camera Camera
	camera.programId = program
	camera.SetModel(mgl32.Translate3D(0, 0, 0))
	camera.SetPosition(mgl32.Vec3{1, 0, 2})
	camera.LookAt(mgl32.Vec3{0, 0, 0})
	camera.SetScreenSize(800, 600)

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		camera.SetScreenSize(float32(width), float32(height))
	})

	prevTime := glfw.GetTime()
	frameCount := 0

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)
		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(planeVertices)/3))

		glfw.PollEvents()
		window.SwapBuffers()

		if window.GetKey(glfw.KeyQ) == glfw.Press {
			break
		}

		currentTime := glfw.GetTime()
		frameCount++

		if currentTime-prevTime >= 1 {
			fmt.Println("Fps: ", frameCount)
			frameCount = 0
			prevTime = currentTime
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
