package main

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
	"time"
)

func GetGrid(filePath string) ([][]color.Color, string, int, int) {
	now := math.Round(float64(time.Now().UnixNano()) / 1000000)
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Could not open the file: ", err)
	}
	defer file.Close()
	openMS := int(math.Round(float64(time.Now().UnixNano())/1000000) - now)
	now2 := math.Round(float64(time.Now().UnixNano()) / 1000000)
	content, format, err := image.Decode(file)
	if err != nil {
		log.Fatal("Could not decode the file: ", err)
	}

	rect := content.Bounds()
	rgba := image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
	draw.Draw(rgba, rgba.Bounds(), content, rect.Min, draw.Src)

	var grid [][]color.Color
	size := rgba.Bounds().Size()
	for i := 0; i < size.X; i += 1 {
		var y []color.Color
		for j := 0; j < size.Y; j += 1 {
			y = append(y, rgba.At(i, j))
		}
		grid = append(grid, y)
	}
	convertMS := int(math.Round(float64(time.Now().UnixNano())/1000000) - now2)
	return grid, format, openMS, convertMS
}
