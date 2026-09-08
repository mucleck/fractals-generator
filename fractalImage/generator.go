package fractalimage

import (
	"flag"
	"image"
	"image/png"
	"math"
	"os"
	"sync"
)

type Config struct {
	Width, Height, Iterations  int
	ViewWidth, Real, Imaginary float64
	C                          complex128
	FileName                   string
	pngCompression             bool

	//Coordinates to map from complex number to a pixel
	xmin, xmax, ymin, ymax float64

	dx float64
	dy float64
}

type Region struct {
	x, y, maxX, maxY int
}

const Threads int = 8

func GenerateImage(config Config) error {
	config.mapBounds()

	myImage := image.NewRGBA(image.Rect(0, 0, config.Width, config.Height))

	var wg sync.WaitGroup

	jobs := make(chan Region)

	for range Threads {
		wg.Go(func() {
			for region := range jobs {
				for i := region.x; i < region.maxX; i++ {
					for j := region.y; j < region.maxY; j++ {
						c := pixelColor(i, j, &config)
						//I do this instead of .Set because we know for sure that the image is inside
						//borders and its already and RGBA
						offset := j*myImage.Stride + i*4

						myImage.Pix[offset] = c
						myImage.Pix[offset+1] = c
						myImage.Pix[offset+2] = c
						myImage.Pix[offset+3] = 255
					}
				}
			}
		})
	}

	regions := config.generateRegions()

	for _, region := range regions {
		jobs <- region
	}

	close(jobs)
	wg.Wait()

	file, err := os.Create(config.FileName)
	if err != nil {
		return err
	}

	defer file.Close()
	var compressionLevel png.CompressionLevel
	if config.pngCompression {
		compressionLevel = png.BestSpeed
	} else {
		compressionLevel = png.NoCompression
	}
	encoder := png.Encoder{
		CompressionLevel: compressionLevel,
	}

	err = encoder.Encode(file, myImage)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) generateRegions() []Region {
	var regions []Region
	for i := range Threads {
		x := i * (c.Width / Threads)
		maxX := (i + 1) * (c.Width / Threads)
		for j := range Threads {
			y, maxY := j*(c.Height/Threads), (j+1)*(c.Height/Threads)
			regions = append(regions, Region{x: x, maxX: maxX, y: y, maxY: maxY})
		}
	}

	return regions
}

func (c *Config) getImaginaryNumber(x, y int) complex128 {
	re := c.xmin + float64(x)*c.dx
	im := c.ymax + float64(y)*c.dy

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
	flag.BoolVar(&c.pngCompression, "compression", false, "Choose if you want to apply png compression")
	flag.Parse()

	c.C = complex(c.Real, c.Imaginary)
}

func pixelColor(x, y int, config *Config) uint8 {
	c := config.getImaginaryNumber(x, y)

	z := complex128(0)
	var i int

	for i = 0; i < config.Iterations && real(z)*real(z)+imag(z)*imag(z) <= 4; i++ {
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

	v := 5 + float64(i) - math.Log2(math.Log(tr+ti))
	// pickColorGrayscale()
	v = 512 * v / float64(config.Iterations)

	if v > 255 {
		v = 255
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
	c.dx = (c.xmax - c.xmin) / float64(c.Width-1)
	c.dy = (c.ymax - c.ymin) / float64(c.Height-1)
}
