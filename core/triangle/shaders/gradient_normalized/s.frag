#version 330 core

uniform vec2 resolution;

out vec4 FragColor;

in vec3 vColor;

void main() {
    // float gradient = vColor.y;
    float gradient = gl_FragCoord.y / resolution.y;
    float normalizedGradient = (gradient + 1.) / 2.;
    FragColor = vec4(normalizedGradient, 0.0, 1.0 - normalizedGradient, 1.0);
}
