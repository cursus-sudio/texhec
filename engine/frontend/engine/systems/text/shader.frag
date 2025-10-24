#version 460 core

in FS {
    vec2 uv;
    flat int glyph;
} fs;

layout(binding = 0) uniform sampler2DArray tex;

out vec4 fragColor;

void main() {
    vec4 color = texture(tex, vec3(fs.uv.xy, fs.glyph));
    if (color.a <= 0.1) {
        discard;
    }
    fragColor = color;
}
