package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.4-core/gl"
	"strings"
)

func compileShader(source string, shaderType uint32) (uint32, error) {
	shaderId := gl.CreateShader(shaderType)

	csource, free := gl.Strs(source)
	gl.ShaderSource(shaderId, 1, csource, nil)
	free()

	gl.CompileShader(shaderId)

	var status int32
	gl.GetShaderiv(shaderId, gl.COMPILE_STATUS, &status)

	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shaderId, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shaderId, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile shader: %v", log)
	}

	return shaderId, nil
}
