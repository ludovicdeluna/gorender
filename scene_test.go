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
)

var msg string = "%s :\nGot:  %v\nWant: %v"

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
				if e := expect(testSize.got); e.Equals(testSize.want) {
					t.Error(e.It("New scene get correct sizes"))
				}
			})
		}
	}
	t.Run("Points", func(t *testing.T) {
		rect := image.Rect(0, 0, size.width, size.height)
		if e := expect(scene.Image.Bounds()); e.Equals(rect) {
			t.Error(e.It("scene.Image as correct rectangle points"))
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
	for point := range pixelIterator(scene.Width) {
		t.Run("pixel_"+point.title, func(t *testing.T) {
			if e := expect(scene.Image.At(point.x, point.y)); e.Equals(testCase) {
				t.Error(e.It("Color for this pixel is random"))
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
	testName = "mustFail"
	t.Run(testName, func(t *testing.T) {
		testCase := testCases[testName]
		if e := expect(scene.Save(testCase)); e.HasError("Can't save file " + testCase) {
			t.Error(e.It("Save on no-writable location will fail"))
		}
	})
	testName = "mustSucceed"
	t.Run(testName, func(t *testing.T) {
		testCase := testCases[testName]
		if e := expect(scene.Save(testCase)); e.Equals(error(nil)) {
			t.Error(e.It("Save must be succeed"))
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

// Expects helper (to avoid repetition -> if got, want := , ; got != want)
type Expects struct {
	got   interface{}
	want  interface{}
	title string
}

func expect(got interface{}) *Expects {
	expects := Expects{got: got}
	return &expects
}

func (e *Expects) String() string {
	return fmt.Sprintf(msg, e.title, e.got, e.want)
}

func (e *Expects) It(m string) string {
	e.title = m
	return e.String()
}

func (e *Expects) Equals(want interface{}) bool {
	e.want = want
	return e.got != e.want
}

func (e *Expects) HasError(want interface{}) bool {
	e.want = want
	switch t := e.got.(type) {
	case error:
		return t.Error() != e.want
	default:
		return true
	}
}
