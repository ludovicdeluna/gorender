package gorender

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ludovicdeluna/gorender/helpers/expect"
)

func TestScene_NewScene(t *testing.T) {
	size := struct{ width, height int }{4, 4}
	scene := NewScene(size.width, size.height)
	{
		testCases := []struct {
			title string
			got   int
			want  int
		}{
			{"Width", scene.Width, size.width},
			{"Height", scene.Height, size.height},
		}
		for _, testSize := range testCases {
			t.Run(testSize.title, func(t *testing.T) {
				if e := expect.For(testSize.got).Equals(testSize.want); e.Fail {
					t.Error(e.It("Initialize new scene with correct height x width"))
				}
			})
		}
	}
	t.Run("Points", func(t *testing.T) {
		rect := image.Rect(0, 0, size.width, size.height)
		if e := expect.For(scene.Image.Bounds()).Equals(rect); e.Fail {
			t.Error(e.It("Have Image object using correct height x width"))
		}
	})
}

func TestScene_EachPixel(t *testing.T) {
	scene := NewScene(4, 4)
	type Point struct {
		x, y  int
		title string
	}
	pixelIterator := func(max int) <-chan Point {
		done := make(chan Point)
		go func() {
			point := Point{0, 0, "0-0"}
			for point.y < max {
				point.title = fmt.Sprintf("%d-%d", point.x, point.y)
				done <- point
				point.x = point.x + 1
				if point.x >= max {
					point.x = 0
					point.y = point.y + 1
				}
			}
			close(done)
		}()
		return done
	}

	testCase := randomColor()
	scene.EachPixel(func(x, y int) color.RGBA {
		return testCase
	})
	e := expect.New()
	for point := range pixelIterator(scene.Width) {
		t.Run("pixel_"+point.title, func(t *testing.T) {
			if e.For(scene.Image.At(point.x, point.y)).Equals(testCase).Fail {
				t.Error(e.It("Colorize all pixels with color function"))
			}
		})
	}
}

func TestScene_Save(t *testing.T) {
	var testName string
	scene := NewScene(4, 4)
	testCases := map[string]string{
		"mustFail":    filepath.Join("cant_write_here", "test_myfile.png"),
		"mustSucceed": os.DevNull,
	}
	e := expect.New()
	testName = "mustFail"
	t.Run(testName, func(t *testing.T) {
		testCase := testCases[testName]
		if e.For(scene.Save(testCase)).HasError("Can't save file " + testCase).Fail {
			t.Error(e.It("Return error when file can't be saved"))
		}
	})
	testName = "mustSucceed"
	t.Run(testName, func(t *testing.T) {
		testCase := testCases[testName]
		if e.For(scene.Save(testCase)).HasNoError().Fail {
			t.Error(e.It("Save the file"))
		}
	})
}

// Helpers
var randomizer = rand.New(rand.NewSource(time.Now().Unix()))

func randomColor() color.RGBA {
	rgb := make([]byte, 3) // Byte are uint8, 8 bits unsigned values -> 0-255
	randomizer.Read(rgb)   // Assigne random 8 bits values into slide (len 3)
	return color.RGBA{rgb[0], rgb[1], rgb[2], byte(255)}
}
