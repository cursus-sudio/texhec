#version 460 core

in FS {
    vec2 uv;
    flat int glyph;
} fs;

layout(binding = 0) uniform sampler2DArray tex;
uniform vec4 u_color;

out vec4 fragColor;

void main() {
    vec4 texColor = texture(tex, vec3(fs.uv.xy, fs.glyph));
    fragColor = texColor * u_color;
}
