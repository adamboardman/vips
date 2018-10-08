package vips

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/image/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"io/ioutil"
	"os"
	"testing"
)

func BenchmarkParallel(b *testing.B) {
	options := Options{Width: 800, Height: 600, Crop: true}
	f, err := os.Open("testdata/1.jpg")
	if err != nil {
		b.Fatal(err)
	}
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := Resize(buf, options)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.StopTimer()
}

func BenchmarkSerialized(b *testing.B) {
	options := Options{Width: 800, Height: 600, Crop: true}
	f, err := os.Open("testdata/1.jpg")
	if err != nil {
		b.Fatal(err)
	}
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Resize(buf, options)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
}

func TestResize(t *testing.T) {
	var testCases = []struct {
		origWidth      int
		origHeight     int
		maxWidth       int
		maxHeight      int
		expectedWidth  uint
		expectedHeight uint
	}{
		{5, 5, 10, 10, 5, 5},
		{10, 10, 5, 5, 5, 5},
		{10, 50, 10, 10, 2, 10},
		{50, 10, 10, 10, 10, 2},
		{50, 100, 60, 90, 45, 90},
		{120, 100, 60, 90, 60, 50},
		{200, 250, 200, 150, 120, 150},
		{200, 250, 120, 0, 120, 150},
		{200, 250, 0, 150, 120, 150},
	}

	for index, mt := range testCases {
		img := image.NewGray16(image.Rect(0, 0, mt.origWidth, mt.origHeight))
		buf := new(bytes.Buffer)
		err := jpeg.Encode(buf, img, nil)
		if err != nil {
			t.Errorf(
				"%d. jpeg.Encode(buf, img, nil) error: %#v",
				index, err)
		}

		options := Options{
			Width:        mt.maxWidth,
			Height:       mt.maxHeight,
			Crop:         false,
			Enlarge:      false,
			Extend:       EXTEND_WHITE,
			Interpolator: NOHALO,
			Gravity:      CENTRE,
			Quality:      90,
		}

		newImg, err := Resize(buf.Bytes(), options)
		if err != nil {
			t.Errorf(
				"%d. Resize(imgData, %#v) error: %#v",
				index, options, err)
		}

		outImg, err := jpeg.Decode(bytes.NewReader(newImg))
		if err != nil {
			t.Errorf(
				"%d. jpeg.Decode(newImg) error: %#v",
				index, err)
		}

		newWidth := uint(outImg.Bounds().Dx())
		newHeight := uint(outImg.Bounds().Dy())
		if newWidth != mt.expectedWidth ||
			newHeight != mt.expectedHeight {
			t.Fatalf("%d. Resize(imgData, %#v) => "+
				"width: %v, height: %v, want width: %v, height: %v, "+
				"originl size: %vx%v",
				index, options,
				newWidth, newHeight, mt.expectedWidth, mt.expectedHeight,
				mt.origWidth, mt.origHeight,
			)
		}
	}
}

func TestResizeGif(t *testing.T) {
	expectedWidth := 100
	expectedHeight := 80
	options := Options{Width: expectedWidth, Height: expectedHeight, Crop: true}
	f, err := os.Open("testdata/8.gif")
	if err != nil {
		t.Fatal(err)
	}
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	origImg, err := gif.Decode(bytes.NewReader(buf))
	if err != nil {
		t.Errorf(
			"GIF. gif.Decode(orgImg) error: %#v",
			err)
	}
	origWidth := int(origImg.Bounds().Dx())
	origHeight := int(origImg.Bounds().Dy())

	newImg, err := Resize(buf, options)
	if err != nil {
		t.Errorf(
			"GIF. Resize(imgData, %#v) error: %#v",
			options, err)
	}

	outImg, err := jpeg.Decode(bytes.NewReader(newImg))
	if err != nil {
		t.Errorf(
			"GIF. jpeg.Decode(newImg) error: %#v",
			err)
	}

	newWidth := int(outImg.Bounds().Dx())
	newHeight := int(outImg.Bounds().Dy())
	if newWidth != expectedWidth ||
		newHeight != expectedHeight {
		t.Fatalf("GIF. Resize(imgData, %#v) => "+
			"width: %v, height: %v, want width: %v, height: %v, "+
			"originl size: %vx%v",
			options,
			newWidth, newHeight, expectedWidth, expectedHeight,
			origWidth, origHeight,
		)
	}
}

func TestResizeWebp(t *testing.T) {
	expectedWidth := 100
	expectedHeight := 80
	options := Options{Width: expectedWidth, Height: expectedHeight, Crop: true}
	f, err := os.Open("testdata/9.webp")
	if err != nil {
		t.Fatal(err)
	}
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	origImg, err := webp.Decode(bytes.NewReader(buf))
	if err != nil {
		t.Errorf(
			"GIF. gif.Decode(orgImg) error: %#v",
			err)
	}
	origWidth := int(origImg.Bounds().Dx())
	origHeight := int(origImg.Bounds().Dy())

	newImg, err := Resize(buf, options)
	if err != nil {
		t.Errorf(
			"GIF. Resize(imgData, %#v) error: %#v",
			options, err)
	}

	outImg, err := jpeg.Decode(bytes.NewReader(newImg))
	if err != nil {
		t.Errorf(
			"GIF. jpeg.Decode(newImg) error: %#v",
			err)
	}

	newWidth := int(outImg.Bounds().Dx())
	newHeight := int(outImg.Bounds().Dy())
	if newWidth != expectedWidth ||
		newHeight != expectedHeight {
		t.Fatalf("GIF. Resize(imgData, %#v) => "+
			"width: %v, height: %v, want width: %v, height: %v, "+
			"originl size: %vx%v",
			options,
			newWidth, newHeight, expectedWidth, expectedHeight,
			origWidth, origHeight,
		)
	}
}

func TestApplyFocusToCropCalc(t *testing.T) {
	Convey("Shrink Horizontal from square", t, func() {
		left, top := applyFocusToCropCalc(1000, 1000, 1000, 200, Focus{0.3, 0.3})
		So(left, ShouldEqual, 0)
		So(top, ShouldEqual, 200)
	})
	Convey("Shrink Vertical from square", t, func() {
		left, top := applyFocusToCropCalc(1000, 1000, 200, 1000, Focus{0.3, 0.3})
		So(left, ShouldEqual, 200)
		So(top, ShouldEqual, 0)
	})
	Convey("Shrink Horizontal from square", t, func() {
		left, top := applyFocusToCropCalc(1000, 1000, 300, 200, Focus{0.3, 0.3})
		So(left, ShouldEqual, 0)
		So(top, ShouldEqual, 0)
	})
	Convey("Shrink Vertical from square", t, func() {
		left, top := applyFocusToCropCalc(1000, 1000, 200, 300, Focus{0.3, 0.3})
		So(left, ShouldEqual, 0)
		So(top, ShouldEqual, 0)
	})
	Convey("Shrink Horizontal from square", t, func() {
		left, top := applyFocusToCropCalc(1000, 1000, 300, 200, Focus{0.7, 0.7})
		So(left, ShouldEqual, 0)
		So(top, ShouldEqual, 334)
	})
	Convey("Shrink Vertical from square", t, func() {
		left, top := applyFocusToCropCalc(1000, 1000, 200, 300, Focus{0.7, 0.7})
		So(left, ShouldEqual, 334)
		So(top, ShouldEqual, 0)
	})
	Convey("Shrink Horizontal", t, func() {
		left, top := applyFocusToCropCalc(1000, 400, 200, 200, Focus{0.3, 0.3})
		So(left, ShouldEqual, 100)
		So(top, ShouldEqual, 0)
	})
	Convey("Shrink Vertical", t, func() {
		left, top := applyFocusToCropCalc(400, 1000, 200, 200, Focus{0.3, 0.3})
		So(left, ShouldEqual, 0)
		So(top, ShouldEqual, 100)
	})
	Convey("Shrink Horizontal", t, func() {
		left, top := applyFocusToCropCalc(1000, 400, 200, 200, Focus{0.7, 0.7})
		So(left, ShouldEqual, 500)
		So(top, ShouldEqual, 0)
	})
	Convey("Shrink Vertical", t, func() {
		left, top := applyFocusToCropCalc(400, 1000, 200, 200, Focus{0.7, 0.7})
		So(left, ShouldEqual, 0)
		So(top, ShouldEqual, 500)
	})
}
