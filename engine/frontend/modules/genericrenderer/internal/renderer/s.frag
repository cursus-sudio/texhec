#version 460 core

in FS {
    vec2 uv;
} fs;

layout(binding = 0) uniform sampler2D tex;
uniform vec4 u_color;

out vec4 fragColor;

void main() {
    vec4 texColor = texture(tex, fs.uv);
    fragColor = texColor * u_color;
}
