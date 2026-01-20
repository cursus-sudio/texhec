#version 460 core

layout(points) in;

layout(triangle_strip, max_vertices = 4) out;

//

in GS {
    flat int vertexID;
} gs_in[];

out FS {
    vec2 uv;
    flat int tileType;
} gs_out;

layout(std430, binding = 0) buffer Grid {
    int grid[];
};

uniform uint width;

vec2 getCoord(int i) {
    return vec2(
        i % width, // X: Coord(index) % g.width
        i / width //  Y: Coord(index) / g.width
    );
}

int getIndex(vec2 coord) {
    return int(coord.x) + int(coord.y) * int(width);
}

//

uniform mat4 mvp;
uniform float widthInv;
uniform float heightInv;

vec2 offsets[4] = vec2[](
        vec2(0.0, 0.0),
        vec2(0.0, 1.0),
        vec2(1.0, 0.0),
        vec2(1.0, 1.0)
    );

void main() {
    int i = gs_in[0].vertexID;

    int tileType = int(grid[i]);
    vec2 pos = getCoord(i);

    gs_out.tileType = tileType;

    for (int i = 0; i < 4; i++) {
        gl_Position = mvp * vec4(
                    widthInv * (pos.x + offsets[i].x) - 1,
                    heightInv * (pos.y + offsets[i].y) - 1,
                    0, 1);
        gs_out.uv = offsets[i];
        EmitVertex();
    }

    EndPrimitive();
}
