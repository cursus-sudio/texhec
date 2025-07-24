#version 330

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vertPos;
in vec2 vertTexCoord;

out vec2 fragTexCoord;

void main() {
    fragTexCoord = vertTexCoord;

    gl_Position = projection * camera * model * vec4(vertPos, 1);
}
