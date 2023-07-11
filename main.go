package main

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"
)

const BORDER int = 5
const RADIUS int = 25
const SAMPLE string = "samples/2.jpg"
const THRESHOLD int = 120

type Point struct {
	IntensityDifference float64
	IsEmpty             bool
	X                   int
	Y                   int
}

func main() {
	img, format := decodeSource(SAMPLE)
	width, height := img.Rect.Max.X, img.Rect.Max.Y

	threshold := uint8(THRESHOLD)

	pixLen := len(img.Pix)
	threads := runtime.NumCPU()
	pixPerThread := getPixPerThread(pixLen, threads)

	var wg sync.WaitGroup

	grayTimeStart := math.Round(float64(time.Now().UnixNano()))
	gray := make([]uint8, len(img.Pix))
	grayscale := func(thread int) {
		defer wg.Done()
		startIndex := pixPerThread * thread
		endIndex := clamp(startIndex+pixPerThread, 0, pixLen)
		for i := startIndex; i < endIndex; i += 4 {
			channel := uint8((int(img.Pix[i]) + int(img.Pix[i+1]) + int(img.Pix[i+2])) / 3)
			gray[i], gray[i+1], gray[i+2], gray[i+3] = channel, channel, channel, img.Pix[i+3]
		}
	}
	for t := 0; t < threads; t += 1 {
		wg.Add(1)
		go grayscale(t)
	}
	wg.Wait()
	fmt.Printf(
		"convert to gray in %f ms\n",
		(math.Round(float64(time.Now().UnixNano()))-grayTimeStart)/1e+6,
	)

	candidatesCount := 0
	points := []Point{}

	fastTimeStart := math.Round(float64(time.Now().UnixNano()))
	fast := func(thread int) {
		defer wg.Done()
		startIndex := pixPerThread * thread
		endIndex := clamp(startIndex+pixPerThread, 0, pixLen)
		for i := startIndex; i < endIndex; i += 4 {
			x, y := getCoordinates(i/4, width)

			// skip border pixels
			if x < BORDER || x > width-BORDER ||
				y < BORDER || y > height-BORDER {
				continue
			}

			circle := [16]uint8{}
			grayPixel := gray[i]
			circle[0] = gray[getPixel(x, y-3, width)]
			circle[4] = gray[getPixel(x+3, y, width)]
			circle[8] = gray[getPixel(x, y+3, width)]
			circle[12] = gray[getPixel(x-3, y, width)]

			deltaMax := uint8(clamp(int(grayPixel)+int(threshold), 0, 255))
			deltaMin := uint8(clamp(int(grayPixel)-int(threshold), 0, 255))

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

			circle[1] = gray[getPixel(x+1, y-3, width)]
			circle[2] = gray[getPixel(x+2, y-2, width)]
			circle[3] = gray[getPixel(x+3, y-1, width)]
			circle[5] = gray[getPixel(x+3, y+1, width)]
			circle[6] = gray[getPixel(x+2, y+2, width)]
			circle[7] = gray[getPixel(x+1, y+3, width)]
			circle[9] = gray[getPixel(x-1, y+3, width)]
			circle[10] = gray[getPixel(x-2, y+2, width)]
			circle[11] = gray[getPixel(x-3, y+1, width)]
			circle[13] = gray[getPixel(x-3, y-1, width)]
			circle[14] = gray[getPixel(x-2, y-2, width)]
			circle[15] = gray[getPixel(x-1, y-3, width)]

			invalidIndexes := make([]int, 0, 12)
			for index, value := range circle {
				if value < deltaMax && value > deltaMin {
					invalidIndexes = append(invalidIndexes, index)
				}
			}

			if len(invalidIndexes) > 1 {
				if len(invalidIndexes) > 4 {
					continue
				}

				checkBright := true
				if darkerCount > brighterCount {
					checkBright = false
				}

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
						intensitySum += math.Abs(float64(grayPixel) - float64(point))
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

				if maxValid < 12 {
					continue
				}

				intensityAverage := intensitySum / float64(maxValid)
				intensityDifference := float64(grayPixel) - intensityAverage
				if intensityAverage > float64(grayPixel) {
					intensityDifference = intensityAverage - float64(grayPixel)
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
	for t := 0; t < threads; t += 1 {
		wg.Add(1)
		go fast(t)
	}
	wg.Wait()
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
		drawSquare(img.Pix, pointsToDrawY[i].X, pointsToDrawY[i].Y, width)
	}

	encodeImage(img, format)
}
