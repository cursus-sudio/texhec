#version 330 core

uniform sampler2D tex;

in int type;

layout(location = 0) in vec4 fragColor;
layout(location = 1) in vec2 fragTexCoord;

out vec4 fragColor;

// void

void main() {
    fragColor = color;
}
