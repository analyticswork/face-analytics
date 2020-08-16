package main

import (
	"fmt"
	"path/filepath"

	"github.com/analyticswork/face-analytics/recognizer"
)

const (
	photosDir        = "photos"
	dataDir          = "models"
	defaultTolerance = 0.5
)

func main() {
	rec := recognizer.Recognizer{
		Tolerance: defaultTolerance,
		UseGray:   true,
		UseCNN:    false,
	}
	if err := rec.Init(dataDir); err != nil {
		fmt.Println(err)
		return
	}
	defer rec.Close()

	addFile(&rec, filepath.Join(photosDir, "amy.jpg"), "Amy")
	addFile(&rec, filepath.Join(photosDir, "bernadette.jpg"), "Bernadette")
	addFile(&rec, filepath.Join(photosDir, "howard.jpg"), "Howard")
	addFile(&rec, filepath.Join(photosDir, "penny.jpg"), "Penny")
	addFile(&rec, filepath.Join(photosDir, "raj.jpg"), "Raj")
	addFile(&rec, filepath.Join(photosDir, "sheldon.jpg"), "Sheldon")
	addFile(&rec, filepath.Join(photosDir, "leonard.jpg"), "Leonard")

	rec.SetSamples()

	faces, err := rec.ClassifyMultiples(filepath.Join(photosDir, "elenco3.jpg"))
	if err != nil {
		fmt.Println(err)
		return
	}

	img, err := rec.DrawFaces(filepath.Join(photosDir, "elenco3.jpg"), faces)
	if err != nil {
		fmt.Println(err)
		return
	}
	rec.SaveImage("faces.jpg", img)
}

func addFile(rec *recognizer.Recognizer, path, id string) {
	if err := rec.AddImageToDataset(path, id); err != nil {
		fmt.Println(err)
		return
	}
}
