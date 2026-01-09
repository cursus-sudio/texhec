#version 460 core

layout(points) in;

layout(triangle_strip, max_vertices = 4) out;

//

in GS {
    flat float x;
    flat float y;
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
    float xIn = gs_in[0].x;
    float yIn = gs_in[0].y;
    int tileType = gs_in[0].tileType;

    // shared outputs
    gs_out.tileType = tileType;

    for (int cornerX = 0; cornerX < 2; cornerX++) {
        for (int cornerY = 0; cornerY < 2; cornerY++) {
            float x = xIn + float(cornerX);
            float y = yIn + float(cornerY);

            vec4 pos = vec4(x * tileSize, y * tileSize, gridDepth, 1.);

            // unique outputs
            gl_Position = camera * pos;
            gs_out.uv = vec2(cornerX, cornerY);
            EmitVertex();
        }
    }
    EndPrimitive();
}
