#version 460 core

in FS {
    flat int id;
    vec2 uv;
} fs;

layout(binding = 0) uniform sampler2DArray tex;

layout(std430, binding = 1) buffer Color {
    vec4 colors[];
};

layout(std430, binding = 2) buffer Frame {
    int frames[];
};

out vec4 fragColor;

void main() {
    vec4 finalColor = texture(tex, vec3(fs.uv, frames[fs.id]));
    finalColor = finalColor * colors[fs.id];

    if (finalColor.a == 0) {
        discard;
    }
    fragColor = finalColor;
}
