package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/cmplx"
	"os"
)

type Mandelbrot struct {
	width  int
	height int
	iters  int

	xmin float64
	xmax float64
	ymin float64
	ymax float64
}

func main() {

	mandelbrot := NewMandelbrot(
		3840, 2160, 100, complex(-0.3, 1.9*2), 3.4*2,
	)

	myImage := image.NewRGBA(image.Rect(0, 0, mandelbrot.width, mandelbrot.height))

	for i := range mandelbrot.width {
		for j := range mandelbrot.height {
			c := mandelbrot.pixelColor(i, j)
			myImage.Set(i, j, color.RGBA{c, c, c, 255})
		}
	}

	file, err := os.Create("mandelbrot.png")
	if err != nil {
		log.Fatalf("Error creating a file: %v", err)
	}

	defer file.Close()
	if err := png.Encode(file, myImage); err != nil {
		log.Fatalf("Can't encode the image data: %v", err)
	}
}

func NewMandelbrot(width, height, iters int, center complex128, viewWidth float64) Mandelbrot {
	ratio := float64(width) / float64(height)

	xmin := real(center) - viewWidth/2
	xmax := real(center) + viewWidth/2

	viewHeight := viewWidth / ratio

	ymin := imag(center) + viewHeight/2
	ymax := imag(center) - viewHeight/2

	return Mandelbrot{
		width:  width,
		height: height,
		iters:  iters,
		xmin:   xmin,
		xmax:   xmax,
		ymin:   ymin,
		ymax:   ymax,
	}
}

func (m Mandelbrot) getImaginaryNumber(x, y int) complex128 {
	re := m.xmin + float64(x)/float64(m.width-1)*(m.xmax-m.xmin)
	im := m.ymax + float64(y)/float64(m.height-1)*(m.ymax-m.ymin)

	return complex(re, im)
}

func (m Mandelbrot) pixelColor(x, y int) uint8 {
	c := m.getImaginaryNumber(x, y)

	z := complex128(0)
	var i int

	for i = 0; i < m.iters && cmplx.Abs(z) <= 2; i++ {
		z = z*z + c
	}

	if i == m.iters {
		return 0
	}

	tr := real(z) * real(z)
	ti := imag(z) * imag(z)

	for range 4 {
		z = z*z + c

	}
	tr = real(z) * real(z)
	ti = imag(z) * imag(z)

	v := 5 +
		float64(i) -
		math.Log(math.Log(tr+ti))/math.Log(2)

	// pickColorGrayscale()
	v = math.Floor(512 * v / float64(m.iters))

	if v > 255 {
		v = 255
	}
	if v < 0 {
		v = 0
	}

	return uint8(v)
}
