#version 460 core

layout(location = 0) in vec2 pos;
layout(location = 1) in int glyph;

out GS {
    vec2 pos;
    flat int glyph;
} fs_out;

void main() {
    fs_out.pos = vec2(pos.x, -pos.y);
    fs_out.glyph = glyph;
}
