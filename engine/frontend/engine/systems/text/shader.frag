#version 460 core

in FS {
    vec2 uv;
    flat int glyph;
} fs;

layout(binding = 0) uniform sampler2DArray tex;

out vec4 fragColor;

void main() {
    fragColor = texture(tex, vec3(fs.uv.x, -fs.uv.y, fs.glyph));
}
