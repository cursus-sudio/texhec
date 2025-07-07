#version 330 core

uniform vec2 resolution;

layout(location = 0) in vec3 pos;
layout(location = 1) in float funny;

out vec3 color;

void main() {
    if (funny == 1.) {
        color = pos.yxz;
    } else {
        color = pos;
    }
    gl_Position = vec4(color.xyz, 1.0);
}
