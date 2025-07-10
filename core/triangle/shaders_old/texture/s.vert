#version 330

// uniform mat4 projection;
// uniform mat4 camera;
// uniform mat4 model;

uniform vec2 resolution;

in vec3 vertPos;
in vec2 vertTexCoord;

out vec2 fragTexCoord;

void main() {
    vec3 normalizedPos = vertPos / vec3(resolution, 1);
    // normalizedPos = (normalizedPos * 2) - 1;
    fragTexCoord = vertTexCoord;
    gl_Position = vec4(normalizedPos, 1);
    // gl_Position = projection * camera * model * vec4(vertPos, 1);
}
