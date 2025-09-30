#version 460 core

layout(location = 0) in vec3 pos;
layout(location = 1) in vec2 uv;

uniform mat4 mvp;

out FS {
    vec2 uv;
} fs;

void main() {
    fs.uv = uv;

    gl_Position = mvp * vec4(pos, 1);
}
