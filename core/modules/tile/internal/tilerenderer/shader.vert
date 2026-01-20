#version 460 core

out GS {
    flat int vertexID;
} fs_out;

void main() {
    fs_out.vertexID = gl_VertexID;
}
