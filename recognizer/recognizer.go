package recognizer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"os"

	goFace "github.com/Kagami/go-face"
)

// Data descriptor of the human face.
type Data struct {
	ID         string
	Descriptor goFace.Descriptor
}

// Face holds coordinates and descriptor of the human face.
type Face struct {
	Data
	Rectangle image.Rectangle
}

// A Recognizer creates face descriptors for provided images and
// classifies them into categories.
type Recognizer struct {
	Tolerance float32
	rec       *goFace.Recognizer
	UseCNN    bool
	UseGray   bool
	Dataset   []Data
}

// Init initialise a recognizer interface.
//
func (r *Recognizer) Init(path string) error {
	r.Tolerance = 0.4
	r.UseCNN = false
	r.UseGray = true

	r.Dataset = []Data{}
	rec, err := goFace.NewRecognizer(path)
	if err != nil {
		return err
	}
	r.rec = rec

	return nil
}

// Close frees resources taken by the Recognizer. Safe to call multiple
// times. Don't use Recognizer after close call.
//
func (r *Recognizer) Close() {
	r.rec.Close()
}

// AddImageToDataset add a sample image to the dataset
//
func (r *Recognizer) AddImageToDataset(path, id string) error {
	file := path
	var err error
	if r.UseGray {
		file, err = r.createTempGrayFile(file, id)
		if err != nil {
			return err
		}
		defer os.Remove(file)
	}

	faces := []goFace.Face{}
	if r.UseCNN {
		faces, err = r.rec.RecognizeFileCNN(file)
	} else {
		faces, err = r.rec.RecognizeFile(file)
	}
	if err != nil {
		return err
	}
	if len(faces) == 0 {
		return errors.New("no face found on the image")
	}
	if len(faces) > 1 {
		return errors.New("No single face on the image")
	}

	f := Data{
		ID:         id,
		Descriptor: faces[0].Descriptor,
	}
	r.Dataset = append(r.Dataset, f)

	return nil
}

// SetSamples sets known descriptors so you can classify the new ones.
//
func (r *Recognizer) SetSamples() {
	samples := []goFace.Descriptor{}
	avengers := []int32{}
	for i, f := range r.Dataset {
		samples = append(samples, f.Descriptor)
		avengers = append(avengers, int32(i))
	}
	r.rec.SetSamples(samples, avengers)
}

// RecognizeSingle returns face if it's the only face on the image or nil otherwise.
// Only JPEG format is currently supported.
//
func (r *Recognizer) RecognizeSingle(path string) (goFace.Face, error) {
	file := path
	var err error
	if r.UseGray {
		file, err = r.createTempGrayFile(file, "64ab59ac42d69274f06eadb11348969e")
		if err != nil {
			return goFace.Face{}, err
		}
		defer os.Remove(file)
	}

	idFace := &goFace.Face{}
	if r.UseCNN {
		idFace, err = r.rec.RecognizeSingleFileCNN(file)
	} else {
		idFace, err = r.rec.RecognizeSingleFile(file)
	}
	if err != nil {
		return goFace.Face{}, fmt.Errorf("can not recognize face err: %v", err)

	}
	if idFace == nil {
		return goFace.Face{}, fmt.Errorf("No single face found on the image")
	}
	return *idFace, nil
}

// RecognizeMultiples returns all faces found on the provided image, sorted from
// left to right. Empty list is returned if there are no faces, error is
// returned if there was some error while decoding/processing image.
// Only JPEG format is currently supported.
//
func (r *Recognizer) RecognizeMultiples(path string) ([]goFace.Face, error) {
	file := path
	var err error
	if r.UseGray {
		file, err = r.createTempGrayFile(file, "64ab59ac42d69274f06eadb11348969e")
		if err != nil {
			return nil, err
		}
		defer os.Remove(file)
	}

	idFaces := []goFace.Face{}
	if r.UseCNN {
		idFaces, err = r.rec.RecognizeFileCNN(file)
	} else {
		idFaces, err = r.rec.RecognizeFile(file)
	}
	if err != nil {
		return nil, fmt.Errorf("can not recognize face. err: %v", err)
	}

	return idFaces, nil

}

//Classify returns all faces identified in the image. Empty list is returned if no match.
//
func (r *Recognizer) Classify(path string) ([]Face, error) {
	face, err := r.RecognizeSingle(path)
	if err != nil {
		return nil, err
	}
	personID := r.rec.ClassifyThreshold(face.Descriptor, r.Tolerance)
	if personID < 0 {
		return nil, fmt.Errorf("can not classify")
	}
	facesRec := []Face{}
	aux := Face{
		Data:      r.Dataset[personID],
		Rectangle: face.Rectangle,
	}
	facesRec = append(facesRec, aux)

	return facesRec, nil
}

// ClassifyMultiples returns all faces identified in the image. Empty list is returned if no match.
//
func (r *Recognizer) ClassifyMultiples(path string) ([]Face, error) {
	faces, err := r.RecognizeMultiples(path)
	if err != nil {
		return nil, fmt.Errorf("can not recognize faces. err: %v", err)
	}

	facesRec := []Face{}
	for _, f := range faces {
		personID := r.rec.ClassifyThreshold(f.Descriptor, r.Tolerance)
		if personID < 0 {
			continue
		}
		aux := Face{
			Data:      r.Dataset[personID],
			Rectangle: f.Rectangle,
		}
		facesRec = append(facesRec, aux)
	}

	return facesRec, nil
}

// fileExists check se file exist
//
func fileExists(fileName string) bool {
	file, err := os.Stat(fileName)
	return (err == nil) && !file.IsDir()
}

// jsonMarshal Marshal interface to array of byte
//
func jsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
