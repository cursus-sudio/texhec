#version 460 core

layout(location = 0) in int x;
layout(location = 1) in int y;
layout(location = 2) in int z;
layout(location = 3) in int tileType;

out GS {
    flat int x;
    flat int y;
    flat int z;
    flat int tileType;
} fs_out;

void main() {
    fs_out.x = x;
    fs_out.y = y;
    fs_out.z = z;
    fs_out.tileType = tileType;
}
