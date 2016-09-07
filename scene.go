package gorender

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

type Scene struct {
	Width  int
	Height int
	Image  *image.RGBA
}

func NewScene(width, height int) *Scene {
	return &Scene{
		Width:  width,
		Height: height,
		Image:  image.NewRGBA(image.Rect(0, 0, width, height)),
	}
}

func (s *Scene) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Can't save file %s", filename)
	}
	defer file.Close()
	return png.Encode(file, s.Image)
}

type ColorFunction func(x, y int) color.RGBA

func (s *Scene) EachPixel(yeld ColorFunction) {
	for x := 0; x < s.Width; x++ {
		for y := 0; y < s.Height; y++ {
			s.Image.Set(x, y, yeld(x, y))
		}
	}
}
