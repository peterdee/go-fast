package main

import (
	"fmt"
	"image/color"
)

func clampMax[T uint8](value, max T) T {
	if value > max {
		return max
	}
	return value
}

func getPartials(pixel color.Color) (uint8, uint8, uint8, uint8) {
	r, g, b, a := pixel.RGBA()
	return uint8(r), uint8(g), uint8(b), uint8(a)
}

func main() {
	grid, format, _, _ := GetGrid("samples/2.jpg")

	height, width := len(grid[0]), len(grid)

	border := 5
	threshold := 84

	gray := make([][]color.Color, width)
	for x := 0; x < width; x += 1 {
		row := make([]color.Color, height)
		for y := 0; y < height; y += 1 {
			r, g, b, a := getPartials(grid[x][y])
			grayColor := uint8(float32(r+g+b) / 3.0)
			row[y] = color.RGBA{grayColor, grayColor, grayColor, a}
		}
		gray[x] = row
	}

	totalCount := 0
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

			deltaMax := clampMax(pixelGray+uint8(threshold), 255)
			deltaMin := pixelGray - uint8(threshold)

			brighterCount, darkerCount := 0, 0
			if b0Gray > deltaMax {
				brighterCount += 1
			} else if b0Gray < deltaMin {
				darkerCount += 1
			}
			if b4Gray > deltaMax {
				brighterCount += 1
			} else if b4Gray < deltaMin {
				darkerCount += 1
			}
			if b8Gray > deltaMax {
				brighterCount += 1
			} else if b8Gray < deltaMin {
				darkerCount += 1
			}
			if b12Gray > deltaMax {
				brighterCount += 1
			} else if b12Gray < deltaMin {
				darkerCount += 1
			}
			if brighterCount != 3 && darkerCount != 3 {
				continue
			}
			totalCount += 1
			fmt.Println(x, y)

			grid[x][y-1] = color.RGBA{255, 0, 0, 255}
			grid[x+1][y-1] = color.RGBA{255, 0, 0, 255}
			grid[x+1][y] = color.RGBA{255, 0, 0, 255}
			grid[x+1][y+1] = color.RGBA{255, 0, 0, 255}
			grid[x][y+1] = color.RGBA{255, 0, 0, 255}
			grid[x-1][y+1] = color.RGBA{255, 0, 0, 255}
			grid[x-1][y] = color.RGBA{255, 0, 0, 255}
			grid[x-1][y-1] = color.RGBA{255, 0, 0, 255}
		}
	}

	fmt.Println("total", totalCount)
	SaveGrid(format, grid)
}
