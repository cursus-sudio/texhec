#version 460 core

layout(points) in;

layout(triangle_strip, max_vertices = 4) out;

//

in GS {
    flat int x;
    flat int y;
    flat int z;
    flat int tileType;
} gs_in[];

out FS {
    vec2 uv;
    flat int tileType;
} gs_out;

//

uniform mat4 camera;
uniform int tileSize;
uniform float gridDepth;

void main() {
    int xIn = gs_in[0].x;
    int yIn = gs_in[0].y;
    int zIn = gs_in[0].z;
    int tileType = gs_in[0].tileType;

    // shared outputs
    gs_out.tileType = tileType;

    for (int cornerX = 0; cornerX < 2; cornerX++) {
        for (int cornerY = 0; cornerY < 2; cornerY++) {
            int x = xIn + cornerX;
            int y = yIn + cornerY;

            vec4 pos = vec4(x * tileSize, y * tileSize, zIn + gridDepth, 1.);

            // unique outputs
            gl_Position = camera * pos;
            gs_out.uv = vec2(cornerX, cornerY);
            EmitVertex();
        }
    }
    EndPrimitive();
}
