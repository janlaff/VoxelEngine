#version 440 core

layout(location = 0) in vec3 vertexPosition;

uniform mat4 view;
uniform mat4 projection;
uniform mat4 model;

out mat4 inverseMVP;

void main() {
    mat4 MVP = projection * view * model;
    inverseMVP = inverse(MVP);
    gl_Position = MVP * vec4(vertexPosition, 1);
}