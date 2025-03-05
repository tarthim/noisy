package noisy

import (
	"fmt"
	"image/color"
	"math/rand"
)

func parseHexColor(s string) (color.RGBA, error) {
	var c color.RGBA
	if !isValidHex(s) {
		return c, fmt.Errorf("%s is not a valid hex code", s)
	}
	_, err := fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	if err != nil {
		return c, err
	}
	c.A = 255
	return c, nil
}

func isValidHex(s string) bool {
	if len(s) != 7 || s[0] != '#' {
		return false
	}
	for _, c := range s[1:] {
		if !((c >= '0' && c <= '9') ||
			(c >= 'A' && c <= 'F') ||
			(c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}

func rgbaToArray(rgba color.RGBA) [4]uint8 {
	return [4]uint8{rgba.R, rgba.G, rgba.B, rgba.A}
}

func randomIntArray8() [4]uint8 {
	return [4]uint8{
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
	}
}
