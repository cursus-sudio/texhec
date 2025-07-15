#version 330

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

uniform vec3 resolution;

in vec3 vertPos;
in vec2 vertTexCoord;

out vec2 fragTexCoord;

void main() {
    fragTexCoord = vertTexCoord;

    vec3 normalizedPos = vertPos / resolution;
    gl_Position = projection * camera * model * vec4(vertPos, 1);
}
