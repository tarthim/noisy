package noisy

import (
	"fmt"
	"image"
	"image/png"
	"math/rand/v2"
	"os"
	"sync"
)

type noiseimagetype []uint8

type noisy struct {
	imageData noiseimagetype
	width     int
	height    int
	image     *image.RGBA
	operation operation
	data      noiseData
	filename  string
}

type noiseData interface{}

type whiteNoiseData struct {
	color1 [4]uint8
	color2 [4]uint8
	chance float64
}
type colorNoiseData struct{}

type simplexNoiseData struct {
	color1 [4]uint8
	color2 [4]uint8
	scale  float64
}

type operation int

const (
	White operation = iota
	Color
	Simplex
)

var (
	operationMap = make(map[string]operation)
)

func init() {
	operationMap["white"] = White
	operationMap["color"] = Color
	operationMap["simplex"] = Simplex
}

func New(width int, height int, c1 string, c2 string, chance float64, operation string, filename string, simplex float64) (*noisy, error) {
	// Extract operation
	op, err := translateOperation(operation)
	if err != nil {
		return nil, err
	}

	n := &noisy{
		width:     width,
		height:    height,
		operation: op,
		filename:  filename,
	}

	ch1, err := parseHexColor(c1)
	if err != nil {
		return nil, err
	}
	ch2, err := parseHexColor(c2)
	if err != nil {
		return nil, err
	}
	cd1 := rgbaToArray(ch1)
	cd2 := rgbaToArray(ch2)

	switch op {
	case White:
		n.data = &whiteNoiseData{
			color1: cd1,
			color2: cd2,
			chance: chance,
		}
	case Color:
		n.data = &colorNoiseData{}
	case Simplex:
		n.data = &simplexNoiseData{
			scale:  simplex,
			color1: cd1,
			color2: cd2,
		}
	}

	err = n.validate()
	if err != nil {
		return nil, err
	}

	n.fillImage()

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// img.Pix is the underlying image in image; this is an []uint8
	img.Pix = n.imageData
	n.image = img

	return n, nil
}

func (n *noisy) SaveAsPNG() {
	file, err := os.Create(n.filename + ".png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := png.Encode(file, n.image); err != nil {
		panic(err)
	}
}

func (n *noisy) GetImage() *image.RGBA {
	return n.image
}

func (n *noisy) fillImage() {
	rows := n.height
	cols := n.width
	// Create a flat slice to hold RGBA data
	// imageData := make([]uint8, rows*cols*4)
	if n.operation == Simplex {
		if simplexNoiseData, ok := n.data.(*simplexNoiseData); ok {
			n.imageData = generateSimplexNoise(cols, rows, simplexNoiseData.scale, simplexNoiseData.color1, simplexNoiseData.color2)
			return
		}
	}

	if n.operation == White || n.operation == Color {
		n.generateWhiteNoise(cols, rows)
		return
	}
}

func (n *noisy) generateWhiteNoise(cols, rows int) {
	imageData := make([]uint8, rows*cols*4)
	var wg sync.WaitGroup
	for i := range rows {
		wg.Add(1)
		go func(row int) {
			defer wg.Done()
			for col := range cols {
				aC := n.getNextColor()
				index := (row*cols + col) * 4
				imageData[index] = aC[0]   // Red
				imageData[index+1] = aC[1] // Green
				imageData[index+2] = aC[2] // Blue
				imageData[index+3] = aC[3] // Alpha
			}
		}(i)
	}
	wg.Wait()
	n.imageData = imageData
}

func (n *noisy) getNextColor() [4]uint8 {
	if n.operation == Color {
		return randomIntArray8()
	}

	if n.operation == White {
		// Type assertion
		if whiteNoiseData, ok := n.data.(*whiteNoiseData); ok {
			c := whiteNoiseData.chance
			if rand.Float64() > c {
				return whiteNoiseData.color2
			}
			return whiteNoiseData.color1
		}
	}
	return [4]uint8{0, 0, 0, 0}
}

func (n *noisy) validate() error {
	switch n.operation {
	case White:
		return n.isValidWhiteNoise()
	case Color:
		return n.isValidColorNoise()
	case Simplex:
		return n.isValidSimplexNoise()
	default:
		return fmt.Errorf("cannot determine operation mode for Noisy")
	}
}

func (n *noisy) isValidWhiteNoise() error {
	err := n.validateDimensions()
	if err != nil {
		return err
	}
	return nil
}

func (n *noisy) isValidColorNoise() error {
	err := n.validateDimensions()
	if err != nil {
		return err
	}
	return nil
}

func (n *noisy) isValidSimplexNoise() error {
	err := n.validateDimensions()
	if err != nil {
		return err
	}
	return nil
}

func (n *noisy) validateDimensions() error {
	if n.height < 0 {
		return fmt.Errorf("height cannot be under 0")
	}
	if n.width < 0 {
		return fmt.Errorf("width cannot be under 0")
	}
	return nil
}

func translateOperation(opStr string) (operation, error) {
	op, exists := operationMap[opStr]
	if exists {
		return op, nil
	}
	return -1, fmt.Errorf("operation is unknown")
}
