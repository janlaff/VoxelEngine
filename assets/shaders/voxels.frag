#version 440 core

out vec4 color;

uniform mat4 mvp;
uniform mat4 inverseViewProjection;
uniform float screenWidth;
uniform float screenHeight;

struct Ray {
    vec3 origin;
    vec3 direction;
};

vec3 walkRay(Ray ray, float distance) {
    return ray.origin + ray.direction * distance;
}

float sphereDistance(vec3 center, float radius, Ray ray) {
    vec3 oc = ray.origin - center;
    float a = dot(ray.direction, ray.direction);
    float b = 2 * dot(oc, ray.direction);
    float c = dot(oc, oc) - radius*radius;
    float discriminant = b*b - 4*a*c;

    if (discriminant < 0) {
        return -1;
    } else {
        return (-b - sqrt(discriminant)) / (2*a);
    }
}

vec3 rayColor(Ray ray) {
    float t = sphereDistance(vec3(0, 0, -1), 0.5, ray);
    if (t > 0) {
        vec3 n = normalize(walkRay(ray, t) - vec3(0, 0, -1));
        return 0.5 * vec3(n.x+1, n.y+1, n.z+1);
    }

    vec3 unitDirection = normalize(ray.direction);
    t = 0.5*(unitDirection.y + 1.0);
    return (1.0-t)*vec3(1, 1, 1)+t*vec3(0.5, 0.7, 1.0);
}

void main() {
    float aspectRatio = 16.0 / 9.0;
    float width = 800;
    float height = 600;

    float viewportHeight = 2.0;
    float viewportWidth = aspectRatio * viewportHeight;
    float focalLength = 1;

    vec3 origin = vec3(0);
    vec3 horizontal = vec3(viewportWidth, 0, 0);
    vec3 vertical = vec3(0, viewportHeight, 0);
    vec3 lowerLeftCorner = origin - horizontal/2 - vertical/2 - vec3(0, 0, focalLength);

    float u = gl_FragCoord.x / (width - 1);
    float v = gl_FragCoord.y / (height - 1);

    Ray eye;
    eye.origin = origin;
    eye.direction = lowerLeftCorner + u*horizontal + v*vertical - origin;

    color = vec4(rayColor(eye), 1.0);
}