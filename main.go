package main

func main() {
	grid, format, openMS, convertMS := GetGrid("samples/1.jpg")

	height, width := len(grid[0]), len(grid)

	border := 5

	for x := border; x < width-border; x += 1 {
		for y := border; y < height-border; y += 1 {
			pixel := grid[x][y]
			b0 := grid[x][y-3]
			b1 := grid[x+1][y-3]
			b2 := grid[x+2][y-2]
			b3 := grid[x+3][y-1]
			b4 := grid[x+3][y]
			b5 := grid[x+3][y+1]
			b6 := grid[x+2][y+2]
			b7 := grid[x+1][y+3]
			b8 := grid[x][y+3]
			b9 := grid[x-1][y+3]
			b10 := grid[x-2][y+2]
			b11 := grid[x-3][y+1]
			b12 := grid[x-3][y]
			b13 := grid[x-3][y-1]
			b14 := grid[x-2][y-2]
			b15 := grid[x-1][y-3]
		}
	}
}
