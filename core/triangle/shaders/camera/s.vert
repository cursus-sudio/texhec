#version 330 core

uniform vec2 resolution;
uniform mat4 camera;

layout(location = 0) in vec3 vertexPos;
layout(location = 1) in vec4 vertexColor;

out vec4 color;

void main() {
    vec3 normalizedPos = vertexPos / vec3(resolution.x, resolution.y, 1);

    normalizedPos = (normalizedPos * 2) - 1;
    normalizedPos *= vec3(1., -1., 1.);

    color = vertexColor;
    gl_Position = vec4(normalizedPos.xyz, 1.0);
}
