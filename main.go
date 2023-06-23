package main

import (
	"fmt"
	"image/color"
)

const BORDER int = 5
const SAMPLE string = "samples/3.png"
const THRESHOLD int = 115

func main() {
	grid, format, _, _ := GetGrid("samples/3.png")

	height, width := len(grid[0]), len(grid)

	border := BORDER
	threshold := uint8(THRESHOLD)

	gray := make([][]color.Color, width)
	for x := 0; x < width; x += 1 {
		row := make([]color.Color, height)
		for y := 0; y < height; y += 1 {
			r, g, b, a := getPartials(grid[x][y])
			grayColor := uint8((float32(r) + float32(g) + float32(b)) / 3.0)
			row[y] = color.RGBA{grayColor, grayColor, grayColor, uint8(a)}
		}
		gray[x] = row
	}

	candidatesCount := 0
	for x := border; x < width-border; x += 1 {
		for y := border; y < height-border; y += 1 {
			pixelGray, _, _, _ := getPartials(gray[x][y])
			b0Gray, _, _, _ := getPartials(gray[x][y-3])
			b4Gray, _, _, _ := getPartials(gray[x+3][y])
			b8Gray, _, _, _ := getPartials(gray[x][y+3])
			b12Gray, _, _, _ := getPartials(gray[x-3][y])
			// b1 := gray[x+1][y-3]
			// b2 := gray[x+2][y-2]
			// b3 := gray[x+3][y-1]
			// b5 := gray[x+3][y+1]
			// b6 := gray[x+2][y+2]
			// b7 := gray[x+1][y+3]
			// b9 := gray[x-1][y+3]
			// b10 := gray[x-2][y+2]
			// b11 := gray[x-3][y+1]
			// b13 := gray[x-3][y-1]
			// b14 := gray[x-2][y-2]
			// b15 := gray[x-1][y-3]

			deltaMax := uint8(clamp(int(pixelGray)+int(threshold), 0, 255))
			deltaMin := uint8(clamp(int(pixelGray)-int(threshold), 0, 255))

			b0Valid, b4Valid, b8Valid, b12Valid := false, false, false, false

			brighterCount, darkerCount := 0, 0
			if b0Gray > deltaMax {
				brighterCount += 1
				b0Valid = true
			} else if b0Gray < deltaMin {
				darkerCount += 1
				b0Valid = true
			}
			if b4Gray > deltaMax {
				brighterCount += 1
				b4Valid = true
			} else if b4Gray < deltaMin {
				darkerCount += 1
				b4Valid = true
			}
			if b8Gray > deltaMax {
				brighterCount += 1
				b8Valid = true
			} else if b8Gray < deltaMin {
				darkerCount += 1
				b8Valid = true
			}
			if b12Gray > deltaMax {
				brighterCount += 1
				b12Valid = true
			} else if b12Gray < deltaMin {
				darkerCount += 1
				b12Valid = true
			}

			// skip pixel if both counts are insufficient
			if brighterCount < 3 && darkerCount < 3 {
				continue
			}

			candidatesCount += 1

			// TODO: determine if point is valid
			if b0Valid && b4Valid && b8Valid {
				b1Gray, _, _, _ := getPartials(gray[x+1][y-3])
				b2Gray, _, _, _ := getPartials(gray[x+2][y-2])
				b3Gray, _, _, _ := getPartials(gray[x+3][y-1])
				b5Gray, _, _, _ := getPartials(gray[x+3][y+1])
				b6Gray, _, _, _ := getPartials(gray[x+2][y+2])
				b7Gray, _, _, _ := getPartials(gray[x+1][y+3])
				if b1Gray < deltaMax && b1Gray > deltaMin {
					continue
				}
				if b2Gray < deltaMax && b2Gray > deltaMin {
					continue
				}
				if b3Gray < deltaMax && b3Gray > deltaMin {
					continue
				}
				if b5Gray < deltaMax && b5Gray > deltaMin {
					continue
				}
				if b6Gray < deltaMax && b6Gray > deltaMin {
					continue
				}
				if b7Gray < deltaMax && b7Gray > deltaMin {
					continue
				}
				drawSquare(grid, x, y)
			}

			if b4Valid && b8Valid && b12Valid {

			}
		}
	}

	fmt.Println("candidates", candidatesCount)
	SaveGrid(format, grid)
}
