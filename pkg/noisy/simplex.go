package noisy

import (
	"math"
	"math/rand"
)

var perm = generatePerm()

func generatePerm() [512]int {
	var perm [512]int
	p := rand.Perm(256)
	for i, v := range p {
		perm[i] = v
		perm[i+256] = v
	}
	return perm
}

func fastfloor(x float64) int {
	if x > 0 {
		return int(x)
	}
	return int(x) - 1
}

var grad3 = [12][2]float64{
	{1, 1}, {-1, 1}, {1, -1}, {-1, -1},
	{1, 0}, {-1, 0}, {1, 0}, {-1, 0},
	{0, 1}, {0, -1}, {0, 1}, {0, -1},
}

func dot(g [2]float64, x, y float64) float64 {
	return g[0]*x + g[1]*y
}

func simplex2D(x, y float64) float64 {
	var F2 = 0.5 * (math.Sqrt(3.0) - 1.0)
	var G2 = (3.0 - math.Sqrt(3.0)) / 6.0

	s := (x + y) * F2
	i := fastfloor(x + s)
	j := fastfloor(y + s)
	t := float64(i+j) * G2
	X0 := float64(i) - t
	Y0 := float64(j) - t
	x0 := x - X0
	y0 := y - Y0

	var i1, j1 int
	if x0 > y0 {
		i1 = 1
		j1 = 0
	} else {
		i1 = 0
		j1 = 1
	}

	x1 := x0 - float64(i1) + G2
	y1 := y0 - float64(j1) + G2
	x2 := x0 - 1.0 + 2.0*G2
	y2 := y0 - 1.0 + 2.0*G2

	ii := i & 255
	jj := j & 255
	gi0 := perm[ii+perm[jj]] % 12
	gi1 := perm[ii+i1+perm[jj+j1]] % 12
	gi2 := perm[ii+1+perm[jj+1]] % 12

	n0, n1, n2 := 0.0, 0.0, 0.0
	t0 := 0.5 - x0*x0 - y0*y0
	if t0 >= 0 {
		t0 *= t0
		n0 = t0 * t0 * dot(grad3[gi0], x0, y0)
	}
	t1 := 0.5 - x1*x1 - y1*y1
	if t1 >= 0 {
		t1 *= t1
		n1 = t1 * t1 * dot(grad3[gi1], x1, y1)
	}
	t2 := 0.5 - x2*x2 - y2*y2
	if t2 >= 0 {
		t2 *= t2
		n2 = t2 * t2 * dot(grad3[gi2], x2, y2)
	}
	return 70.0 * (n0 + n1 + n2)
}

func generateSimplexNoise(width, height int, scale float64, bg, fg [4]uint8) []uint8 {
	octaves := 12             // Number of noise layers
	persistence := 0.5        // Amplitude reduction for each subsequent octave
	baseScale := scale / 1000 // Base scale value for frequency

	imageData := make([]uint8, width*height*4)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var noiseSum float64
			amplitude := 1.0
			frequency := 1.0
			maxAmplitude := 0.0

			for o := 0; o < octaves; o++ {
				noiseSum += simplex2D(float64(x)*baseScale*frequency, float64(y)*baseScale*frequency) * amplitude
				maxAmplitude += amplitude
				amplitude *= persistence
				frequency *= 2
			}

			// Normalize the noise value to [-1, 1] then to [0, 1]
			noiseValue := noiseSum / maxAmplitude
			factor := (noiseValue + 1) / 2

			// Linearly interpolate between bg and fg for each channel
			r := uint8(float64(bg[0])*(1-factor) + float64(fg[0])*factor)
			g := uint8(float64(bg[1])*(1-factor) + float64(fg[1])*factor)
			b := uint8(float64(bg[2])*(1-factor) + float64(fg[2])*factor)
			a := uint8(float64(bg[3])*(1-factor) + float64(fg[3])*factor)

			index := (y*width + x) * 4
			imageData[index] = r
			imageData[index+1] = g
			imageData[index+2] = b
			imageData[index+3] = a
		}
	}
	return imageData
}
