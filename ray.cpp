#include <iostream>

#include "color.h"
#include "vec3.h"

int main()
{
    // Image
    int image_width = 2000;
    int image_height = 2000;
    
    FILE * pFile;
    pFile = fopen ("image.ppm", "wb");

    std::fprintf(pFile, "P6\n%d %d\n255\n", image_width, image_height);

    for (int j = 0; j < image_height; j++)
    {
        std::clog << "\rScanlines remaining: " << (image_height - j) << ' ' << std::flush;
        
        for (int i = 0; i < image_width; i++)
        {
            auto pixel_color = color(double(i)/(image_width-1),
            double(j)/(image_height-1),
            1.0);
            auto r = pixel_color.x();
            auto g = pixel_color.y();
            auto b = pixel_color.z();

            unsigned char pixel[3];
            pixel[0] = int(255.999 * r);
            pixel[1] = int(255.999 * g);
            pixel[2] = int(255.999 * b);
            
            fwrite (pixel , 1, 3, pFile);
        }
    }

    
    fclose (pFile);

    std::clog << "\rDone. opening file, nerd.                 \n";

    std::system("open image.ppm"); 
}