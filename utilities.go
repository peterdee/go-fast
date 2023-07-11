package main

import (
	"math"
)

func clamp[T float64 | int | uint | uint8](value, min, max T) T {
	if value > max {
		return max
	}
	if value < min {
		return min
	}
	return value
}

func drawSquare(pixels []uint8, x, y, width int) {
	for i := -2; i <= 2; i += 1 {
		px1 := getPixel(x+i, y-2, width)
		px2 := getPixel(x+i, y+2, width)
		pixels[px1], pixels[px1+1], pixels[px1+2] = 255, 0, 0
		pixels[px2], pixels[px2+1], pixels[px2+2] = 255, 0, 0
	}
	for i := -1; i <= 1; i += 1 {
		px1 := getPixel(x-2, y+i, width)
		px2 := getPixel(x+2, y+i, width)
		pixels[px1], pixels[px1+1], pixels[px1+2] = 255, 0, 0
		pixels[px2], pixels[px2+1], pixels[px2+2] = 255, 0, 0
	}
}

func getCoordinates(pixel, width int) (int, int) {
	return pixel % width, int(math.Floor(float64(pixel) / float64(width)))
}

func getPixel(x, y, width int) int {
	return ((y * width) + x) * 4
}

func getPixPerThread(pixLen, threads int) int {
	pixPerThreadRaw := float64(pixLen) / float64(threads)
	module := math.Mod(pixPerThreadRaw, 4.0)
	if module == 0 {
		return int(pixPerThreadRaw)
	}
	return int(pixPerThreadRaw + (float64(threads) - math.Mod(pixPerThreadRaw, 4.0)))
}
