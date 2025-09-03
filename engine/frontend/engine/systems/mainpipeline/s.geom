#version 460 core

layout(triangles) in;

layout(triangle_strip, max_vertices = 30) out;

//

in FS {
    vec2 uv;
    flat int drawID;
} gs_in[];

out FS {
    vec2 uv;
    flat int drawID;
} gs_out;

//

layout(std430, binding = 3) buffer ModelProjectionData {
    int modelProjections[];
};

layout(std430, binding = 4) buffer OrthoData {
    mat4 orthoCameras[];
};

layout(std430, binding = 5) buffer PerspectiveData {
    mat4 perspectiveCameras[];
};

void main() {
    int modelProjection = modelProjections[gs_in[0].drawID];

    if (modelProjection == 0) {
        for (int cameraIndex = 0; cameraIndex < orthoCameras.length(); cameraIndex++) {
            mat4 camera = orthoCameras[cameraIndex];
            for (int i = 0; i < gl_in.length(); i++) {
                gl_Position = camera * gl_in[i].gl_Position;

                gs_out.uv = gs_in[i].uv;
                gs_out.drawID = gs_in[i].drawID;
                EmitVertex();
            }
            EndPrimitive();
        }
    } else {
        for (int cameraIndex = 0; cameraIndex < perspectiveCameras.length(); cameraIndex++) {
            mat4 camera = perspectiveCameras[cameraIndex];
            for (int i = 0; i < gl_in.length(); i++) {
                gl_Position = camera * gl_in[i].gl_Position;

                gs_out.uv = gs_in[i].uv;
                gs_out.drawID = gs_in[i].drawID;
                EmitVertex();
            }
            EndPrimitive();
        }
    }
}
