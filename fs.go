package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"time"
)

func decodeSource(path string) (*image.RGBA, string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Could not open the file: ", err)
	}
	defer file.Close()
	content, format, decodingError := image.Decode(file)
	if decodingError != nil {
		log.Fatal("Could not decode the file: ", err)
	}
	rect := content.Bounds()
	img := image.NewRGBA(rect)
	draw.Draw(img, img.Bounds(), content, rect.Min, draw.Src)
	return img, format
}

func encodeImage(img *image.RGBA, format string) {
	name := fmt.Sprintf(`file-%d.%s`, time.Now().Unix(), format)
	newFile, err := os.Create("results/" + name)
	if err != nil {
		log.Fatal("Could not save the file!")
	}
	defer newFile.Close()
	if format == "png" {
		png.Encode(newFile, img.SubImage(img.Rect))
	} else {
		jpeg.Encode(
			newFile,
			img.SubImage(img.Rect),
			&jpeg.Options{
				Quality: 100,
			},
		)
	}
}
