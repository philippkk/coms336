#include <iostream>
#include <vector>

#include "utils/utils.h"
#include "utils/camera.h"
#include "utils/color.h"
#include "utils/ray.h"
#include "utils/vec3.h"
#include "objects/hittable.h"
#include "objects/hittable_list.h"
#include "objects/sphere.h"

#define IMAGE_WIDTH 800

#define R 0
#define G 1
#define B 2

int main()
{
    hittable_list world;

    world.add(make_shared<sphere>(point3(0,0,-1), 0.5));
    world.add(make_shared<sphere>(point3(0,-100.5,-1), 100));

    camera cam;

    cam.aspect_ratio = 16.0 / 9.0;
    cam.image_width  = IMAGE_WIDTH;

    cam.render(world);    
}

