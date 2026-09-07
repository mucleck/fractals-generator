package main

import (
	"log"
	"mandelbrot.com/mandelbrot/fractalImage"
)

func main() {
	config := fractalimage.Config{}
	config.LoadConfig()

	if err := fractalimage.GenerateImage(config); err != nil {
		log.Fatal(err)
	}

}
