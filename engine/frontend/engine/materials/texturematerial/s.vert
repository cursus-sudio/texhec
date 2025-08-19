#version 460 core

layout(location = 0) in vec3 pos;
layout(location = 1) in vec2 textureCoord;

layout(std430, binding = 2) buffer ModelData {
    mat4 model[];
};

layout(std430, binding = 3) buffer ModelProjectionData {
    int modelProjection[];
};

layout(std430, binding = 4) buffer ProjectionData {
    mat4 projections[];
};

out VS {
    vec2 textureCoord;
    flat int drawID;
} vs;

void main() {
    int id = gl_BaseInstance + gl_InstanceID;
    vs.drawID = id;

    //

    vs.textureCoord = textureCoord;

    //

    int modelProj = modelProjection[id];
    mat4 proj = projections[modelProj];

    mat4 M = model[id];
    vec4 wpos = M * vec4(pos, 1.0);

    gl_Position = proj * wpos;
}
