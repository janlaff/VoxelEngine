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

func calcModelViewProjection(cameraPos mgl32.Vec3, width float32, height float32) mgl32.Mat4 {
	projection := mgl32.Perspective(mgl32.DegToRad(45), width/height, 0.1, 100.0)
	view := mgl32.LookAt(
		cameraPos.X(), cameraPos.Y(), cameraPos.Z(),
		0, 0, 0,
		0, 1, 0,
	)
	model := mgl32.Ident4()
	mvp := projection.Mul4(view).Mul4(model)

	return mvp
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

	cameraPos := mgl32.Vec3{2, 2, 2}
	mvp := calcModelViewProjection(cameraPos, 800, 600)

	mvpLoc := gl.GetUniformLocation(program, gl.Str("mvp\x00"))
	cameraPosLoc := gl.GetUniformLocation(program, gl.Str("cameraPos\x00"))
	screenWidthLoc := gl.GetUniformLocation(program, gl.Str("screenWidth\x00"))
	screenHeightLoc := gl.GetUniformLocation(program, gl.Str("screenHeight\x00"))

	gl.UseProgram(program)
	gl.UniformMatrix4fv(mvpLoc, 1, false, &mvp[0])
	gl.Uniform3fv(cameraPosLoc, 1, &cameraPos[0])
	gl.Uniform1f(screenWidthLoc, 800)
	gl.Uniform1f(screenHeightLoc, 600)

	gl.ClearColor(0.1, 0.1, 0.8, 1.0)

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
		mvp = calcModelViewProjection(cameraPos, float32(width), float32(height))
		gl.UniformMatrix4fv(mvpLoc, 1, false, &mvp[0])
		gl.Uniform1f(screenWidthLoc, float32(width))
		gl.Uniform1f(screenHeightLoc, float32(height))
	})

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
