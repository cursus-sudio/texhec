#version 330

uniform mat4 mvp;

in vec3 vertPos;
in vec2 vertTexCoord;

out vec2 fragTexCoord;

void main() {
    fragTexCoord = vertTexCoord;

    gl_Position = mvp * vec4(vertPos, 1);
}
