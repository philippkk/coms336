#ifndef COLOR_H
#define COLOR_H

#include "vec3.h"
#include "interval.h"
using color = vec3;

void write_color(std::vector<int> &pixel, const color& pixel_color) {
    auto r = pixel_color.x();
    auto g = pixel_color.y();
    auto b = pixel_color.z();

   // Translate the [0,1] component values to the byte range [0,255].
    static const interval intensity(0.000, 0.999);
    int rbyte = int(256 * intensity.clamp(r));
    int gbyte = int(256 * intensity.clamp(g));
    int bbyte = int(256 * intensity.clamp(b));

    // Write out the pixel color components.
    pixel.emplace_back(int(255.999 * r));
    pixel.emplace_back(int(255.999 * g));
    pixel.emplace_back(int(255.999 * b));
}

#endif