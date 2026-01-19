#version 460 core

layout(points) in;

layout(triangle_strip, max_vertices = 4) out;

//

in GS {
    flat vec2 pos;
    flat uint tileType;
} gs_in[];

out FS {
    vec2 uv;
    flat int tileType;
} gs_out;

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
    int tileType = int(gs_in[0].tileType);

    gs_out.tileType = tileType;

    for (int i = 0; i < 4; i++) {
        gl_Position = mvp * vec4(
                    widthInv * (gs_in[0].pos.x + offsets[i].x) - 1,
                    heightInv * (gs_in[0].pos.y + offsets[i].y) - 1,
                    0.0, 1.0);
        gs_out.uv = offsets[i];
        EmitVertex();
    }

    EndPrimitive();
}
