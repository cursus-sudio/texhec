#version 460 core

layout(location = 0) in vec3 pos;
layout(location = 1) in vec2 textureCoord;

layout(std430, binding = 2) buffer ModelData {
    mat4 models[];
};

layout(std430, binding = 3) buffer ModelProjectionData {
    int modelProjections[];
};

layout(std430, binding = 4) buffer ProjectionData {
    mat4 viewProjections[2];
};

out FS {
    vec2 textureCoord;
    flat int drawID;
} fs;

void main() {
    // int id = gl_BaseInstance + gl_InstanceID;
    int id = gl_DrawID;
    fs.drawID = id;

    //

    fs.textureCoord = textureCoord;

    //

    int modelProjection = modelProjections[id];

    mat4 viewProjection = viewProjections[modelProjection];

    mat4 model = models[id];
    // vec4 modelPos = model * vec4(pos, 1);

    gl_Position = viewProjection * model * vec4(pos, 1);
}
