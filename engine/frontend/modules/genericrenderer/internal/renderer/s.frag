#version 460 core

in FS {
    vec2 uv;
} fs;

layout(binding = 0) uniform sampler2D tex;

out vec4 fragColor;

void main() {
    fragColor = texture(tex, fs.uv);
}
