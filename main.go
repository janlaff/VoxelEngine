package main

import (
	_ "embed"
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

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

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

	//https://learnopengl.com/Getting-started/Camera
	/*cameraPosition := mgl32.Vec3{0, 0, 3}
	cameraTarget := mgl32.Vec3{0, 0, 0}
	cameraDirection := cameraPosition.Sub(cameraTarget).Normalize()
	up := mgl32.Vec3{0, 1, 0}
	cameraRight := up.Cross(cameraDirection).Normalize()
	cameraUp := cameraDirection.Cross(cameraRight)*/

	projection := mgl32.Perspective(45.0, 800/600, 0.1, 100.0)
	view := mgl32.LookAt(
		1, 1, 2,
		0, 0, 0,
		0, 1, 0,
	)
	model := mgl32.Ident4()
	mvp := projection.Mul4(view).Mul4(model)

	mvpLocation := gl.GetUniformLocation(program, gl.Str("mvp\x00"))

	gl.ClearColor(0.1, 0.1, 0.8, 1.0)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)
		gl.UniformMatrix4fv(mvpLocation, 1, false, &mvp[0])
		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(planeVertices)/3))

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
