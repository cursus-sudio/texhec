#version 460 core

layout(points) in;

layout(triangle_strip, max_vertices = 4) out;

//

in GS {
    flat int vertexID;
} gs_in[];

out FS {
    vec2 uv;
    flat int textures[4];
} gs_out;

layout(std430, binding = 0) buffer Grid {
    int grid[];
};

layout(std430, binding = 1) buffer TileTextures {
    int tileTextures[];
};
layout(std430, binding = 2) buffer TileTexturesSize {
    int tileTexturesSize[];
};

uniform uint width;
uniform uint height;

vec2 getCoord(int i) {
    return vec2( // this uses +1 because of dual grid system
        i % (width + 1), // X: Coord(index) % g.width
        i / (width + 1) //  Y: Coord(index) / g.width
    );
}

int getIndex(vec2 coord) { // this doesn't use +1 because tile map is still normal
    coord = clamp(coord, vec2(0, 0), vec2(width - 1, height - 1));
    return int(coord.x) + int(coord.y) * int(width);
}

int getTile(vec2 coord) {
    return grid[getIndex(coord)];
}
int getTile(int index) {
    return getTile(getCoord(index));
}

int getTexture(int seed, int i) {
    return tileTextures[i] + seed % tileTexturesSize[i];
}

int seed(int num) {
    uint seed = num;
    seed = seed * 1664525u + 1013904223u;
    seed ^= seed >> 16;
    seed *= 0x85ebca6bu;
    seed ^= seed >> 13;
    seed *= 0xc2b2ae35u;
    seed ^= seed >> 16;
    return int(seed);
}

//

void sort4(inout int a[4]) {
    int tmp;
    #define SWAP(i, j) if (a[i] < a[j]) { tmp = a[i]; a[i] = a[j]; a[j] = tmp; }
    SWAP(0, 1);
    SWAP(2, 3);
    SWAP(0, 2);
    SWAP(1, 3);
    SWAP(1, 2);
    #undef SWAP
}

void setBioms(int index, vec2 coord) {
    int neighbours[4] = {
            getTile(coord + vec2(-1, -1)),
            getTile(coord + vec2(0, -1)),
            getTile(coord + vec2(-1, 0)),
            getTile(coord + vec2(0, 0))
        };
    int bioms[4] = neighbours;
    sort4(bioms);

    int seed = seed(index);
    for (int i = 0; i < 4; i++) {
        int base = bioms[i] * 15;
        int mask = 0;
        for (int n = 0; n < 4; n++) {
            mask |= int(neighbours[n] == bioms[i]) * (1 << n);
        }
        gs_out.textures[i] = getTexture(seed, base + mask);
    }
}

uniform mat4 mvp;
uniform float widthInv; // = 2 / width
uniform float heightInv; // = 2 / height

vec2 corners[4] = vec2[](
        vec2(-.5, -.5),
        vec2(-.5, 0.5),
        vec2(0.5, -.5),
        vec2(0.5, 0.5)
    );

void main() {
    int i = gs_in[0].vertexID;

    vec2 coord = getCoord(i);
    setBioms(i, coord);

    for (int i = 0; i < 4; i++) {
        vec4 pos = vec4(0, 0, 0, 1);
        vec2 offset = vec2(
                coord.x + corners[i].x,
                coord.y + corners[i].y
            );
        pos.x = widthInv * offset.x - 1;
        pos.y = heightInv * offset.y - 1;
        gl_Position = mvp * pos;
        gs_out.uv = corners[i] + vec2(.5, .5);
        EmitVertex();
    }

    EndPrimitive();
}
