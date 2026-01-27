#version 460 core

in FS {
    vec2 uv;
    flat int tileType[4];
} fs;

layout(binding = 0) uniform sampler2DArray tex;

out vec4 fragColor;

void main() {
    for (int i = 0; i < 4; i++) {
        int type = fs.tileType[i];
        fragColor = texture(tex, vec3(fs.uv, type));
        if (fragColor.a != 0) {
            break;
        }
    }
}
