package materials

import "github.com/philippkk/coms336/raytracer/internal/utils"

type Texture interface {
	Value(u, v float64, p utils.Vec3) utils.Vec3
}


