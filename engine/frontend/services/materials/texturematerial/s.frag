#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;

out vec4 color;

void main() {
    color = texture(tex, fragTexCoord);
    // color = vec4(1, 1, 0, 1);
}
