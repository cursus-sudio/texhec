#version 460 core

layout(location = 0) in uint tileType;

out GS {
    flat vec2 pos;
    flat uint tileType;
} fs_out;

uniform int width;

void main() {
    int i = gl_VertexID;

    fs_out.pos = vec2(
            uint(i % width), // X: Coord(index) % g.width
            uint(i / width)); // Y: Coord(index) / g.width
    fs_out.tileType = tileType;
    gl_Position = vec4(0.0, 0.0, 0.0, 1.0);
}
