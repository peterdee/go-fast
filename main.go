package main

import (
	"fmt"
)

const BORDER int = 5
const SAMPLE string = "samples/7.png"
const THRESHOLD int = 40

func main() {
	grid, format, _, _ := GetGrid(SAMPLE)

	height, width := len(grid[0]), len(grid)

	border := BORDER
	threshold := uint8(THRESHOLD)

	gray := make([][]uint8, width)
	for x := 0; x < width; x += 1 {
		row := make([]uint8, height)
		for y := 0; y < height; y += 1 {
			r, g, b, _ := getPartials(grid[x][y])
			grayColor := uint8((float32(r) + float32(g) + float32(b)) / 3.0)
			row[y] = grayColor
		}
		gray[x] = row
	}

	candidatesCount := 0
	pointsCount := 0

	for x := border; x < width-border; x += 1 {
		for y := border; y < height-border; y += 1 {
			pixelGray := gray[x][y]

			circle := [16]uint8{}
			circle[0] = gray[x][y-3]  // 0
			circle[4] = gray[x+3][y]  // 4
			circle[8] = gray[x][y+3]  // 8
			circle[12] = gray[x-3][y] // 12

			deltaMax := uint8(clamp(int(pixelGray)+int(threshold), 0, 255))
			deltaMin := uint8(clamp(int(pixelGray)-int(threshold), 0, 255))

			brighterCount, darkerCount := 0, 0
			if circle[0] > deltaMax {
				brighterCount += 1
			} else if circle[0] < deltaMin {
				darkerCount += 1
			}
			if circle[4] > deltaMax {
				brighterCount += 1
			} else if circle[4] < deltaMin {
				darkerCount += 1
			}
			if circle[8] > deltaMax {
				brighterCount += 1
			} else if circle[8] < deltaMin {
				darkerCount += 1
			}
			if circle[12] > deltaMax {
				brighterCount += 1
			} else if circle[12] < deltaMin {
				darkerCount += 1
			}

			// skip pixel if both counts are insufficient
			if brighterCount < 3 && darkerCount < 3 {
				continue
			}

			candidatesCount += 1

			circle[1] = gray[x+1][y-3]
			circle[2] = gray[x+2][y-2]
			circle[3] = gray[x+3][y-1]
			circle[5] = gray[x+3][y+1]
			circle[6] = gray[x+2][y+2]
			circle[7] = gray[x+1][y+3]
			circle[9] = gray[x-1][y+3]
			circle[10] = gray[x-2][y+2]
			circle[11] = gray[x-3][y+1]
			circle[13] = gray[x-3][y-1]
			circle[14] = gray[x-2][y-2]
			circle[15] = gray[x-1][y-3]

			// find indexes of invalid surrounding points in the circle
			invalidIndexes := make([]int, 0, 12)
			for index, value := range circle {
				if value < deltaMax && value > deltaMin {
					invalidIndexes = append(invalidIndexes, index)
				}
			}

			if len(invalidIndexes) > 1 {
				// skip if there are more than 4 invalid indexes
				if len(invalidIndexes) > 4 {
					continue
				}

				// circleIndex := invalidIndexes[0]
				// if circleIndex == 15 {

				// }
				// for i := 0; i < 16; i += 1 {
				// 	if value < deltaMax && value > deltaMin {
				// 		invalidIndexes = append(invalidIndexes, index)
				// 	}
				// }
			}

			pointsCount += 1
			drawSquare(grid, x, y)
		}
	}

	fmt.Println("candidates", candidatesCount, "points", pointsCount)
	SaveGrid(format, grid)
}
