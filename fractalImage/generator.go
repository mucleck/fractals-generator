package fractalimage

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"
)

type Config struct {
	Width, Height, Iterations  int
	ViewWidth, Real, Imaginary float64
	C                          complex128
	FileName                   string

	//Coordinates to map from complex number to a pixel
	xmin, xmax, ymin, ymax float64
}

func GenerateImage(config Config) error {
	config.mapBounds()

	myImage := image.NewRGBA(image.Rect(0, 0, config.Width, config.Height))

	for i := range config.Width {
		for j := range config.Height {
			c := pixelColor(i, j, &config)
			myImage.Set(i, j, color.RGBA{c, c, c, 255})
		}
	}

	file, err := os.Create(config.FileName)
	if err != nil {
		return err
	}

	defer file.Close()
	if err := png.Encode(file, myImage); err != nil {
		return err
	}

	return nil
}

func (c *Config) getImaginaryNumber(x, y int) complex128 {
	re := c.xmin + float64(x)/float64(c.Width-1)*(c.xmax-c.xmin)
	im := c.ymax + float64(y)/float64(c.Height-1)*(c.ymax-c.ymin)

	return complex(re, im)
}

func (c *Config) LoadConfig() {
	flag.IntVar(&c.Width, "w", 3840, "Set de width of the image")
	flag.IntVar(&c.Height, "h", 2160, "Set de height of the image")
	flag.IntVar(&c.Iterations, "I", 100, "Set the number of iterations for the image")
	flag.Float64Var(&c.Real, "r", -0.3, "Set the real part of the c we will be using")
	flag.Float64Var(&c.Imaginary, "i", 3.8, "Set the imaginary part of the c we will be using")
	flag.Float64Var(&c.ViewWidth, "vW", 6.8, "Define how much of the y axis we will see (divide by two)")
	flag.StringVar(&c.FileName, "n", "mandelbrot.png", "Set the name of the generated image")
	flag.Parse()

	c.C = complex(c.Real, c.Imaginary)
}

func pixelColor(x, y int, config *Config) uint8 {
	c := config.getImaginaryNumber(x, y)

	z := complex128(0)
	var i int

	for i = 0; i < config.Iterations && cmplx.Abs(z) <= 2; i++ {
		z = z*z + c
	}

	if i == config.Iterations {
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
	v = math.Floor(512 * v / float64(config.Iterations))

	if v > 255 {
		v = 255
	}
	if v < 0 {
		v = 0
	}

	return uint8(v)
}

func (c *Config) mapBounds() {
	ratio := float64(c.Width) / float64(c.Height)

	c.xmin = c.Real - c.ViewWidth/2
	c.xmax = c.Real + c.ViewWidth/2

	viewHeight := c.ViewWidth / ratio

	c.ymin = c.Imaginary + viewHeight/2
	c.ymax = c.Imaginary - viewHeight/2
}
