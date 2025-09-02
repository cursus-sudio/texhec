#version 460 core

layout(location = 0) in vec3 pos;
layout(location = 1) in vec2 textureCoord;

layout(std430, binding = 2) buffer ModelData {
    mat4 modelsMatrix[];
};

out FS {
    vec2 textureCoord;
    flat int drawID;
} vs;

void main() {
    // int id = gl_BaseInstance + gl_InstanceID;
    int id = gl_DrawID;
    vs.drawID = id;

    //

    vs.textureCoord = textureCoord;

    mat4 modelMatrix = modelsMatrix[id];

    gl_Position = modelMatrix * vec4(pos, 1);
}
