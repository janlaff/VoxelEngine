package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.4-core/gl"
	"strings"
)

func linkProgram(shaderIds []uint32) (uint32, error) {
	programId := gl.CreateProgram()

	for _, shaderId := range shaderIds {
		gl.AttachShader(programId, shaderId)
	}

	gl.LinkProgram(programId)

	var status int32
	gl.GetProgramiv(programId, gl.LINK_STATUS, &status)

	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(programId, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(programId, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile shader: %v", log)
	}

	for _, shaderId := range shaderIds {
		gl.DetachShader(programId, shaderId)
		gl.DeleteShader(shaderId)
	}

	return programId, nil
}
