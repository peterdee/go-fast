package main

import "image/color"

func clamp[T float64 | int | uint | uint8](value, min, max T) T {
	if value > max {
		return max
	}
	if value < min {
		return min
	}
	return value
}

func drawSquare(grid [][]color.Color, x, y int) {
	red := color.RGBA{255, 0, 0, 255}
	for i := -2; i <= 2; i += 1 {
		grid[x+i][y-2], grid[x+i][y+2] = red, red
	}
	for i := -1; i <= 1; i += 1 {
		grid[x-2][y+i], grid[x+2][y+i] = red, red
	}
}

func getPartials(pixel color.Color) (uint8, uint8, uint8, uint8) {
	r, g, b, a := pixel.RGBA()
	return uint8(r), uint8(g), uint8(b), uint8(a)
}
