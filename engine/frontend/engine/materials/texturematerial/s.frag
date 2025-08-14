#version 460 core

// layout(std430, binding = 1) buffer TexLayerData {
buffer TexLayerData {
    int layer[];
};

uniform sampler2D tex;
// uniform sampler2DArray texs;

in VS {
    vec2 texCoord;
    flat int id;
} fs;

out vec4 color;

void main() {
    color = texture(tex, fs.texCoord);

    // vec3 texCoord = vec3(fs.texCoord, layer[fs.id]);
    // color = texture(texs, texCoord);

    // color = vec4(1, 1, 0, 1);
}
