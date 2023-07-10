package main

import (
	"fmt"
	"math"
	"time"
)

const BORDER int = 5
const RADIUS int = 25
const SAMPLE string = "samples/3.png"
const THRESHOLD int = 80

type Point struct {
	IntensityDifference float64 `json:"intensity"`
	IsEmpty             bool    `json:"-"`
	X                   int     `json:"x"`
	Y                   int     `json:"y"`
}

func main() {
	grid, format, _, _ := GetGrid(SAMPLE)

	height, width := len(grid[0]), len(grid)

	border := BORDER
	threshold := uint8(THRESHOLD)

	grayTimeStart := math.Round(float64(time.Now().UnixNano()))
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
	fmt.Printf(
		"convert to gray in %f ms\n",
		(math.Round(float64(time.Now().UnixNano()))-grayTimeStart)/1e+6,
	)

	candidatesCount := 0
	points := []Point{}

	fastTimeStart := math.Round(float64(time.Now().UnixNano()))
	for x := border; x < width-border; x += 1 {
		for y := border; y < height-border; y += 1 {
			pixelGray := gray[x][y]

			circle := [16]uint8{}
			circle[0] = gray[x][y-3]
			circle[4] = gray[x+3][y]
			circle[8] = gray[x][y+3]
			circle[12] = gray[x-3][y]

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

				checkBright := true
				if darkerCount > brighterCount {
					checkBright = false
				}

				// count continuous valid pixels in a circle
				startIndex := invalidIndexes[0]
				nextIndex := startIndex + 1
				if nextIndex > 15 {
					nextIndex = 0
				}
				currentValid := 0
				maxValid := 0
				intensitySum := 0.0
				for i := 0; i < 15; i += 1 {
					point := circle[nextIndex]
					if (checkBright && point > deltaMax) || (!checkBright && point < deltaMin) {
						currentValid += 1
						intensitySum += math.Abs(float64(pixelGray) - float64(point))
					} else {
						currentValid = 0
					}
					if currentValid > maxValid {
						maxValid = currentValid
					}
					nextIndex += 1
					if nextIndex > 15 {
						nextIndex = 0
					}
				}

				// skip if count is less than 12
				if maxValid < 12 {
					continue
				}

				// get average intensity
				intensityAverage := intensitySum / float64(maxValid)
				intensityDifference := float64(circle[0]) - intensityAverage
				if intensityAverage > float64(circle[0]) {
					intensityDifference = intensityAverage - float64(circle[0])
				}

				points = append(
					points,
					Point{
						IntensityDifference: intensityDifference,
						IsEmpty:             false,
						X:                   x,
						Y:                   y,
					},
				)
			}
		}
	}

	fmt.Printf(
		"find all candidates in %f ms\n",
		(math.Round(float64(time.Now().UnixNano()))-fastTimeStart)/1e+6,
	)

	nmsXTimeStart := math.Round(float64(time.Now().UnixNano()))
	pointsToDrawX := nms(
		points,
		RADIUS,
		Point{
			IsEmpty: true,
		},
		[]Point{},
		[][]Point{},
		false,
		'x',
	)
	nmsXTime := (math.Round(float64(time.Now().UnixNano())) - nmsXTimeStart) / 1e+6
	nmsYTimeStart := math.Round(float64(time.Now().UnixNano()))
	pointsToDrawY := nms(
		pointsToDrawX,
		RADIUS,
		Point{
			IsEmpty: true,
		},
		[]Point{},
		[][]Point{},
		false,
		'y',
	)
	nmsYTime := (math.Round(float64(time.Now().UnixNano())) - nmsYTimeStart) / 1e+6

	fmt.Println(
		"candidates:",
		candidatesCount,
		"\npoints before NMS:",
		len(points),
		"\npoints after NMS:",
		len(pointsToDrawY),
	)

	fmt.Printf(
		"NMS time: %f ms (x) + %f ms (y) = %f ms (total)\n",
		nmsXTime,
		nmsYTime,
		nmsXTime+nmsYTime,
	)

	for i := range pointsToDrawY {
		drawSquare(grid, pointsToDrawY[i].X, pointsToDrawY[i].Y)
	}

	SaveGrid(format, grid)
}
