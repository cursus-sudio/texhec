#version 460 core

layout(location = 0) in float x;
layout(location = 1) in float y;
layout(location = 2) in int tileType;

out GS {
    flat float x;
    flat float y;
    flat int tileType;
} fs_out;

void main() {
    fs_out.x = x;
    fs_out.y = y;
    fs_out.tileType = tileType;
}
