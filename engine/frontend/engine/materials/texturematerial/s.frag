#version 460 core

in VS {
    vec2 textureCoord;
    flat int drawID;
} fs;

layout(std430, binding = 1) buffer ModelTexData {
    int texturesIDs[];
};

layout(binding = 0) uniform sampler2DArray texs;

out vec4 fragColor;

void main() {
    int textureID = texturesIDs[fs.drawID];
    vec3 base = texture(texs, vec3(fs.textureCoord, float(textureID))).rgb;
    fragColor = vec4(base, 1.0);
}
