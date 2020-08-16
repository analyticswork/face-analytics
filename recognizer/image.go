package recognizer

import (
	"crypto/rand"
	"encoding/hex"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"

	goFace "github.com/Kagami/go-face"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

// LoadImage load an image from file
//
func (r *Recognizer) LoadImage(path string) (image.Image, error) {
	existingImageFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer existingImageFile.Close()

	imageData, _, err := image.Decode(existingImageFile)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}

// SaveImage save an image to jpeg file
//
func (r *Recognizer) SaveImage(path string, Img image.Image) error {
	outputFile, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := jpeg.Encode(outputFile, Img, nil); err != nil {
		return err
	}
	return outputFile.Close()
}

// GrayScale Convert an image to grayscale
//
func (r *Recognizer) GrayScale(imgSrc image.Image) image.Image {
	return imaging.Grayscale(imgSrc)
}

// createTempGrayFile create a temporary image in grayscale
//
func (r *Recognizer) createTempGrayFile(path, id string) (string, error) {
	name := r.tempFileName(id, ".jpeg")
	img, err := r.LoadImage(path)
	if err != nil {
		return "", err
	}
	img = r.GrayScale(img)
	err = r.SaveImage(name, img)
	if err != nil {
		return "", err
	}
	return name, nil
}

// tempFileName generates a temporary filename
//
func (r *Recognizer) tempFileName(prefix, suffix string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix)
}

// DrawFaces draws the faces identified in the original image
//
func (r *Recognizer) DrawFaces(path string, faces []Face) (image.Image, error) {
	img, err := r.LoadImage(path)
	if err != nil {
		return nil, err
	}
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(font, &truetype.Options{Size: 24})

	dc := gg.NewContextForImage(img)
	dc.SetFontFace(face)

	for _, f := range faces {
		dc.SetRGB255(0, 0, 255)

		x := float64(f.Rectangle.Min.X)
		y := float64(f.Rectangle.Min.Y)
		w := float64(f.Rectangle.Dx())
		h := float64(f.Rectangle.Dy())

		dc.DrawString(f.Id, x, y+h+20)

		dc.DrawRectangle(x, y, w, h)
		dc.SetLineWidth(4.0)
		dc.SetStrokeStyle(gg.NewSolidPattern(color.RGBA{R: 0, G: 0, B: 255, A: 255}))
		dc.Stroke()

	}
	img = dc.Image()
	return img, nil
}

// DrawFaces2 draws the faces in the original image
//
func (r *Recognizer) DrawFaces2(path string, faces []goFace.Face) (image.Image, error) {
	aux := make([]Face, 0)
	for _, f := range faces {
		auxFace := Face{}
		auxFace.Rectangle = f.Rectangle
		auxFace.Descriptor = f.Descriptor
		aux = append(aux, auxFace)
	}

	return r.DrawFaces(path, aux)
}
