package main

import "github.com/go-gl/mathgl/mgl32"

const ChunkSize = 16

type Chunk struct {
	data     [ChunkSize * ChunkSize * ChunkSize]bool
	position mgl32.Vec3
}

/*func (c *Chunk) InRange(position mgl32.Vec3) bool {
	if position.X() < c.position.X() {
		return true
	}
	if position.X() + ChunkSize > c.position.X() {
		return true
	}
}

func (c *Chunk) IsSolid(position mgl32.Vec3) bool {
	if !InRange(position) {
		return true
	}

	index := int(position.Z())*ChunkSize*ChunkSize + int(position.Y()) + int(position.X())
	return c.data[index]
	return true
}*/

func (c *Chunk) CreateMesh() MeshData {
	var meshData MeshData

	for z := 0; z < ChunkSize; z++ {
		for y := 0; y < ChunkSize; y++ {
			for x := 0; x < ChunkSize; x++ {
				index := z*ChunkSize*ChunkSize + y*ChunkSize + x
				position := mgl32.Vec3{float32(x), float32(y), float32(z)}

				if c.data[index] {
					meshData.AddCube(position)
				}
			}
		}
	}

	return meshData
}
