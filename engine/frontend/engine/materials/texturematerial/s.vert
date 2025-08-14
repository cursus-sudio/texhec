#version 460 core

uniform mat4 mvp;

in vec3 pos;
in vec2 texCoord;

out VS {
    vec2 texCoord;
    flat int id;
} fs;

void main() {
    int id = gl_DrawID;
    fs.id = id;

    fs.texCoord = texCoord;

    gl_Position = mvp * vec4(pos, 1);
}
