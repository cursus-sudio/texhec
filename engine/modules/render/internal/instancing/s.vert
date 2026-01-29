#version 460 core

layout(location = 0) in vec3 pos;
layout(location = 1) in vec2 uv;

layout(std430, binding = 0) buffer Models {
    mat4 models[];
};

layout(std430, binding = 3) buffer Groups {
    uint groups[];
};

uniform mat4 camera;
uniform uint cameraGroups;

out FS {
    flat int id;
    vec2 uv;
} fs;

void main() {
    int id = gl_InstanceID;

    uint groups = groups[id];
    if ((groups & cameraGroups) == 0) {
        gl_Position = vec4(2, 2, 2, 1);
        return;
    }
    vec4 pos = camera * models[id] * vec4(pos, 1);

    fs.id = id;
    fs.uv = uv;
    gl_Position = pos;
}
