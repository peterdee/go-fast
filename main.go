package main

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"
)

const BORDER int = 5
const RADIUS int = 15
const SAMPLE string = "samples/3.jpg"
const SAVE_GRAYSCALE bool = true
const THRESHOLD uint8 = 90
const USE_NMS bool = true

type Point struct {
	IntensitySum float64
	IsEmpty      bool
	X            int
	Y            int
}

func main() {
	img, format := decodeSource(SAMPLE)
	width, height := img.Rect.Max.X, img.Rect.Max.Y

	pixLen := len(img.Pix)
	threads := runtime.NumCPU()
	pixPerThread := getPixPerThread(pixLen, threads)

	var mu = &sync.Mutex{}
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
		"Grayscale time: %f ms\n",
		(math.Round(float64(time.Now().UnixNano()))-grayTimeStart)/1e+6,
	)

	points := []Point{}

	fastTimeStart := math.Round(float64(time.Now().UnixNano()))
	fast := func(thread int) {
		defer wg.Done()
		startIndex := pixPerThread * thread
		endIndex := clamp(startIndex+pixPerThread, 0, pixLen)
		for i := startIndex; i < endIndex; i += 4 {
			x, y := getCoordinates(i/4, width)

			// TODO: handle border pixels, skip for now
			if x < BORDER || x > width-BORDER ||
				y < BORDER || y > height-BORDER {
				continue
			}

			// high-speed test
			brighterCount, darkerCount := 0, 0
			grayPixel := gray[i]
			deltaMax := uint8(clamp(int(grayPixel)+int(THRESHOLD), 0, 255))
			deltaMin := uint8(clamp(int(grayPixel)-int(THRESHOLD), 0, 255))
			point0, point8 := gray[getPixel(x, y-3, width)], gray[getPixel(x, y+3, width)]
			if point0 > deltaMax {
				brighterCount += 1
			} else if point0 < deltaMin {
				darkerCount += 1
			}
			if point8 > deltaMax {
				brighterCount += 1
			} else if point8 < deltaMin {
				darkerCount += 1
			}
			if brighterCount+darkerCount == 0 {
				continue
			}
			point4, point12 := gray[getPixel(x+3, y, width)], gray[getPixel(x-3, y, width)]
			if point4 > deltaMax {
				brighterCount += 1
			} else if point4 < deltaMin {
				darkerCount += 1
			}
			if point12 > deltaMax {
				brighterCount += 1
			} else if point12 < deltaMin {
				darkerCount += 1
			}
			if brighterCount < 3 && darkerCount < 3 {
				continue
			}

			circle := [16]uint8{}
			circle[0] = point0
			circle[1] = gray[getPixel(x+1, y-3, width)]
			circle[2] = gray[getPixel(x+2, y-2, width)]
			circle[3] = gray[getPixel(x+3, y-1, width)]
			circle[4] = point4
			circle[5] = gray[getPixel(x+3, y+1, width)]
			circle[6] = gray[getPixel(x+2, y+2, width)]
			circle[7] = gray[getPixel(x+1, y+3, width)]
			circle[8] = point8
			circle[9] = gray[getPixel(x-1, y+3, width)]
			circle[10] = gray[getPixel(x-2, y+2, width)]
			circle[11] = gray[getPixel(x-3, y+1, width)]
			circle[12] = point12
			circle[13] = gray[getPixel(x-3, y-1, width)]
			circle[14] = gray[getPixel(x-2, y-2, width)]
			circle[15] = gray[getPixel(x-1, y-3, width)]

			invalidIndexes := make([]int, 0, 12)
			for index, value := range circle {
				if value < deltaMax && value > deltaMin {
					invalidIndexes = append(invalidIndexes, index)
				}
			}

			invalidIndexesLength := len(invalidIndexes)
			if invalidIndexesLength > 4 {
				continue
			}

			checkBright := darkerCount < brighterCount
			nextIndex := 0
			if invalidIndexesLength > 0 {
				nextIndex = clamp(invalidIndexes[0]+1, 0, 15)
			}
			currentValid := 0
			maxValid := 0
			intensitySum := 0.0
			for i := 0; i < 15; i += 1 {
				point := circle[nextIndex]
				if (checkBright && point > deltaMax) || (!checkBright && point < deltaMin) {
					currentValid += 1
				} else {
					currentValid = 0
				}
				if currentValid > maxValid {
					maxValid = currentValid
				}
				intensitySum += math.Abs(float64(grayPixel) - float64(point))
				nextIndex = clamp(nextIndex+1, 0, 15)
			}

			if maxValid < 12 {
				continue
			}

			mu.Lock()
			points = append(
				points,
				Point{
					IntensitySum: intensitySum,
					IsEmpty:      false,
					X:            x,
					Y:            y,
				},
			)
			mu.Unlock()
		}
	}
	for t := 0; t < threads; t += 1 {
		wg.Add(1)
		go fast(t)
	}
	wg.Wait()
	fmt.Printf(
		"FAST time: %f ms\n",
		(math.Round(float64(time.Now().UnixNano()))-fastTimeStart)/1e+6,
	)

	fmt.Printf(
		"Points: %d (before NMS)\n",
		len(points),
	)

	drawPoints := points

	if USE_NMS {
		nmsTimeStart := math.Round(float64(time.Now().UnixNano()))
		nmsPoints := nmsRecursion(points, RADIUS, 0, true)
		nmsTime := (math.Round(float64(time.Now().UnixNano())) - nmsTimeStart) / 1e+6
		fmt.Printf(
			"NMS time: %f ms\n",
			nmsTime,
		)
		fmt.Printf(
			"Points: %d (after NMS)\n",
			len(nmsPoints),
		)
		drawPoints = nmsPoints
	}

	fmt.Printf(
		"Total processing time: %f ms\n",
		(math.Round(float64(time.Now().UnixNano()))-grayTimeStart)/1e+6,
	)

	if SAVE_GRAYSCALE {
		img.Pix = gray
	}

	for i := range drawPoints {
		drawSquare(img.Pix, drawPoints[i].X, drawPoints[i].Y, width)
	}

	encodeImage(img, format)
}
