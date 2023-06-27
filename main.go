package main

import (
	"fmt"
	"sort"
)

const BORDER int = 5
const DISTANCE int = 10
const SAMPLE string = "samples/3.png"
const THRESHOLD int = 110

type Point struct {
	IntensityDifference float64
	X, Y                int
}

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
	points := []Point{}

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

				checkBright := true
				if darkerCount > brighterCount {
					checkBright = false
				}

				// TODO: fix intensity calculation
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
					intensitySum += float64(point)
					if (checkBright && point > deltaMax) || (!checkBright && point < deltaMin) {
						currentValid += 1
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
				intensityAverage := intensitySum / 15
				intensityDifference := float64(circle[0]) - intensityAverage
				if intensityAverage > float64(circle[0]) {
					intensityDifference = intensityAverage - float64(circle[0])
				}

				points = append(
					points,
					Point{
						IntensityDifference: intensityDifference,
						X:                   x,
						Y:                   y,
					},
				)
			}
		}
	}

	pc := [][]Point{}
	// ci := 0
	for i := range points {
		if len(pc) == 0 {
			pc = append(pc, []Point{points[i]})
			continue
		}
		// TODO: better clustering
		// lcp := pc[ci][len(pc[ci])-1]
		// for j := i; i < len(points); i += 1 {

		// }
	}
	// sort points by coordinates
	// sort.Slice(points, func(i, j int) bool {
	// 	return points[i].X > points[j].X && points[i].Y > points[j].Y
	// })

	// fix point clustering
	pointClusters := [][]Point{}
	clusterIndex := 0
	for _, point := range points {
		if len(pointClusters) == 0 {
			pointClusters = append(pointClusters, []Point{point})
			continue
		}
		lastClusterPoint := pointClusters[clusterIndex][len(pointClusters[clusterIndex])-1]
		maxX, minX := lastClusterPoint.X, point.X
		if point.X > lastClusterPoint.X {
			maxX, minX = point.X, lastClusterPoint.X
		}
		maxY, minY := lastClusterPoint.Y, point.Y
		if point.Y > lastClusterPoint.Y {
			maxY, minY = point.Y, lastClusterPoint.Y
		}
		if maxX-minX < DISTANCE && maxY-minY < DISTANCE {
			pointClusters[clusterIndex] = append(pointClusters[clusterIndex], point)
		} else {
			clusterIndex += 1
			pointClusters = append(pointClusters, []Point{point})
		}
	}

	// sort items in the clusters
	for _, cluster := range pointClusters {
		if len(cluster) == 1 {
			continue
		}
		sort.Slice(cluster, func(i, j int) bool {
			return cluster[i].IntensityDifference > cluster[j].IntensityDifference
		})
	}

	fmt.Println("candidates", candidatesCount, "points", len(points))
	fmt.Println(len(pointClusters), pointClusters)

	for _, cluster := range pointClusters {
		drawSquare(grid, cluster[0].X, cluster[0].Y)
	}
	SaveGrid(format, grid)
}
