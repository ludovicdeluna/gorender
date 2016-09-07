package gorender

import (
	"image"
	"image/color"
	"math/rand"
	"os"
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
				if testSize.got != testSize.want {
					t.Errorf(msg, "New scene get correct sizes", testSize.got, testSize.want)
				}
			})
		}
	}
	t.Run("Points", func(t *testing.T) {
		rect := image.Rect(0, 0, size.width, size.height)
		if got, want := scene.Image.Bounds(), rect; got != want {
			t.Errorf(msg, "scene.Image as correct rectangle points", got, want)
		}
	})
}

func TestScene_EachPixel(t *testing.T) {
	scene := NewScene(4, 4)
	testCase := randomColor()
	scene.EachPixel(func(x, y int) color.RGBA {
		return testCase
	})
	// Be aware: At() use index position starting at 0 (and not 1)
	if got, want := scene.Image.At(3, 3), testCase; got != want {
		t.Errorf(msg, "Random color for pixel 4:4 is random", got, want)
	}
}

func TestScene_Save(t *testing.T) {
	var testName string
	scene := NewScene(4, 4)
	testCases := map[string]string{
		"mustFail":    "cant_write_here/test_myfile.png",
		"mustSucceed": "testdata/test_myfile.png",
	}
	testName = "mustFail"
	t.Run(testName, func(t *testing.T) {
		testCase := testCases[testName]
		want := "Can't save file " + testCase
		got := scene.Save(testCase)
		switch {
		case got == nil:
			t.Error("Didn't failed")
		case got.Error() != want:
			t.Errorf(msg, "Save on no-writable location will fail", got, want)
		}
	})
	testName = "mustSucceed"
	t.Run(testName, func(t *testing.T) {
		testCase := testCases[testName]
		if got, want := scene.Save(testCase), error(nil); got != want {
			t.Errorf(msg, "Save must be succeed (if folder exists)", got, want)
		}
		os.Remove(testCase)
	})
}

// Helpers
var randomizer = rand.New(rand.NewSource(time.Now().Unix()))

func randomColor() color.RGBA {
	rgb := make([]byte, 3) // Byte are uint8, 8 bits unsigned values -> 0-255
	randomizer.Read(rgb)   // Assigne random 8 bits values into slide (len 3)
	return color.RGBA{rgb[0], rgb[1], rgb[2], byte(255)}
}
