#include <iostream>

#include "color.h"
#include "vec3.h"

#define IMAGE_WIDTH 2000
#define IMAGE_HEIGHT 2000

#define R 0
#define G 1
#define B 2

int main()
{
    FILE * pFile;
    pFile = fopen ("image.ppm", "wb");

    std::fprintf(pFile, "P6\n%d %d\n255\n", IMAGE_WIDTH, IMAGE_HEIGHT);

    for (int j = 0; j < IMAGE_HEIGHT; j++)
    {
        std::clog << "\rScanlines remaining: " << (IMAGE_HEIGHT - j) << ' ' << std::flush;
        for (int i = 0; i < IMAGE_WIDTH; i++)
        {
            auto pixel_color = color(double(i)/(IMAGE_WIDTH-1),
            double(j)/(IMAGE_HEIGHT-1),
            1.0);

            auto r = pixel_color.x();
            auto g = pixel_color.y();
            auto b = pixel_color.z();

            unsigned char pixel[3];
            pixel[R] = int(255.999 * r);
            pixel[G] = int(255.999 * g);
            pixel[B] = int(255.999 * b);
            
            fwrite (pixel , 1, 3, pFile);
        }
    }

    
    fclose (pFile);

    std::clog << "\rDone. opening file, nerd.                 \n";

    std::system("open image.ppm"); 
}