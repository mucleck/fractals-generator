package fractalimage

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"testing"
)

func TestGenerateImage(t *testing.T) {
	config := Config{}
	config.LoadConfig()
	config.FileName = "testing.png"
	config.pngCompression = true

	if err := GenerateImage(config); err != nil {
		t.Errorf("err: %v", err)
	}

	got, err := getHash(config.FileName)
	if err != nil {
		t.Errorf("err: %v", err)
	}

	want, err := getHash("golden_image.png")
	if err != nil {
		t.Errorf("err: %v", err)
	}

	if got != want {
		t.Errorf("The golden_image does not match the version of the program:\ngot: %v\nwant: %v", got, want)
	}

	if err := os.Remove(config.FileName); err != nil {
		t.Errorf("Couldnt remove the file generated bby this test but the program is correct: %v", err)
	}

}

// Dummy bench test, rn 190ms per image without compression
func BenchmarkGenerateImage(b *testing.B) {
	config := Config{
		Width:          3840,
		Height:         2160,
		Iterations:     100,
		Real:           -0.3,
		Imaginary:      3.8,
		ViewWidth:      6.8,
		FileName:       "testing.png",
		pngCompression: false,
	}

	config.C = complex(config.Real, config.Imaginary)

	for b.Loop() {
		if err := GenerateImage(config); err != nil {
			b.Fatal(err)
		}
	}
}

func getHash(s string) (string, error) {
	file, err := os.Open(s)
	if err != nil {
		return "", err
	}

	hash := md5.New()

	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
