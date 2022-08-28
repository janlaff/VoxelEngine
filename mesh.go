package main

import "github.com/go-gl/mathgl/mgl32"

type MeshData struct {
	vertices  []float32
	triangles []uint32
}

func (m *MeshData) AddCube(position mgl32.Vec3) {
	m.AddFrontFace(position)
	m.AddRightFace(position)
	m.AddBackFace(position)
	m.AddLeftFace(position)
	m.AddTopFace(position)
	m.AddBottomFace(position)
}

func (m *MeshData) AddVertex(vertex mgl32.Vec3) {
	m.vertices = append(m.vertices, vertex.X(), vertex.Y(), vertex.Z())
}

func (m *MeshData) AddQuadTriangles() {
	m.triangles = append(
		m.triangles,
		uint32(len(m.vertices)/3-4),
		uint32(len(m.vertices)/3-3),
		uint32(len(m.vertices)/3-2),
		uint32(len(m.vertices)/3-4),
		uint32(len(m.vertices)/3-2),
		uint32(len(m.vertices)/3-1),
	)
}

func (m *MeshData) AddFrontFace(position mgl32.Vec3) {
	m.AddVertex(position.Add(mgl32.Vec3{-0.5, 0.5, 0.5}))
	m.AddVertex(position.Add(mgl32.Vec3{0.5, 0.5, 0.5}))
	m.AddVertex(position.Add(mgl32.Vec3{0.5, -0.5, 0.5}))
	m.AddVertex(position.Add(mgl32.Vec3{-0.5, -0.5, 0.5}))
	m.AddQuadTriangles()
}

func (m *MeshData) AddRightFace(position mgl32.Vec3) {
	m.AddVertex(position.Add(mgl32.Vec3{0.5, 0.5, 0.5}))
	m.AddVertex(position.Add(mgl32.Vec3{0.5, 0.5, -0.5}))
	m.AddVertex(position.Add(mgl32.Vec3{0.5, -0.5, -0.5}))
	m.AddVertex(position.Add(mgl32.Vec3{0.5, -0.5, 0.5}))
	m.AddQuadTriangles()
}

func (m *MeshData) AddBackFace(position mgl32.Vec3) {
	m.AddVertex(position.Add(mgl32.Vec3{0.5, 0.5, -0.5}))
	m.AddVertex(position.Add(mgl32.Vec3{-0.5, 0.5, -0.5}))
	m.AddVertex(position.Add(mgl32.Vec3{-0.5, -0.5, -0.5}))
	m.AddVertex(position.Add(mgl32.Vec3{0.5, -0.5, -0.5}))
	m.AddQuadTriangles()
}

func (m *MeshData) AddLeftFace(position mgl32.Vec3) {
	m.AddVertex(position.Add(mgl32.Vec3{-0.5, 0.5, -0.5}))
	m.AddVertex(position.Add(mgl32.Vec3{-0.5, 0.5, 0.5}))
	m.AddVertex(position.Add(mgl32.Vec3{-0.5, -0.5, 0.5}))
	m.AddVertex(position.Add(mgl32.Vec3{-0.5, -0.5, -0.5}))
	m.AddQuadTriangles()
}

func (m *MeshData) AddTopFace(position mgl32.Vec3) {
	m.AddVertex(position.Add(mgl32.Vec3{-0.5, 0.5, -0.5}))
	m.AddVertex(position.Add(mgl32.Vec3{0.5, 0.5, -0.5}))
	m.AddVertex(position.Add(mgl32.Vec3{0.5, 0.5, 0.5}))
	m.AddVertex(position.Add(mgl32.Vec3{-0.5, 0.5, 0.5}))
	m.AddQuadTriangles()
}

func (m *MeshData) AddBottomFace(position mgl32.Vec3) {
	m.AddVertex(position.Add(mgl32.Vec3{-0.5, -0.5, 0.5}))
	m.AddVertex(position.Add(mgl32.Vec3{0.5, -0.5, 0.5}))
	m.AddVertex(position.Add(mgl32.Vec3{0.5, -0.5, -0.5}))
	m.AddVertex(position.Add(mgl32.Vec3{-0.5, -0.5, -0.5}))
	m.AddQuadTriangles()
}
