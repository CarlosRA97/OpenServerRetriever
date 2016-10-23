package worker

import (
	"io/ioutil"
	"strings"
)

type Image struct {
	ID     string `bson:"_id"`
	Name   string `bson:"name"`
	Format string `bson:"format"`
	Data   []byte `bson:"data"`
}

func NewImage(name string, data []byte, format string) Image {
	img := Image{
		ID:     name,
		Name:   name,
		Data:   data,
		Format: format,
	}
	return img
}

func ImageUploader(filepath string) {

	filePath := strings.Split(filepath, "/")
	file := strings.Split(filePath[1], ".")
	fileName := file[0]
	fileExtension := file[1]

	image, err := ioutil.ReadFile(filepath)
	Check(err)

	encodedImage := base64Encode(image)
	img := NewImage(fileName, encodedImage, fileExtension)
	PutData(img, "images")
}
