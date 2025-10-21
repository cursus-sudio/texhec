#version 460 core

layout(points) in;

layout(triangle_strip, max_vertices = 4) out;

//

in GS {
    vec2 pos;
    flat int glyph;
} gs_in[];

out FS {
    vec2 uv;
    flat int glyph;
} gs_out;

uniform mat4 mvp;
uniform vec2 offset;

layout(std430, binding = 0) buffer GlyphWidths {
    float glyphWidths[];
};

//

void main() {
    vec3 posIn = vec3(gs_in[0].pos, 0.);
    int glyph = gs_in[0].glyph;

    vec2 size = vec2(max(glyphWidths[glyph], 1), 1.);

    // shared output
    gs_out.glyph = glyph;

    for (int cornerX = 0; cornerX < 2; cornerX++) {
        for (int cornerY = 0; cornerY < 2; cornerY++) {
            vec4 pos = vec4(
                    size.x * (posIn.x + cornerX) + offset.x,
                    size.y * (posIn.y + cornerY) + offset.y,
                    posIn.z,
                    1.);

            // unique output
            gs_out.uv = vec2(cornerX, cornerY);
            gl_Position = mvp * pos;
            EmitVertex();
        }
    }
    EndPrimitive();
}
