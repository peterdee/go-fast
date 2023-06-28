package main

import (
	"math"
	"sort"
)

// Non-max suppression (recursive)
func nmsr(radius int, initial, result []Point) []Point {
	if len(initial) == 0 {
		return result
	}
	point, rest := initial[0], initial[1:]
	cluster, leftovers := []Point{point}, []Point{}
	for _, element := range rest {
		if int(math.Abs(float64(element.X)-float64(point.X))) < radius &&
			int(math.Abs(float64(element.Y)-float64(point.Y))) < radius {
			cluster = append(cluster, element)
		} else {
			leftovers = append(leftovers, element)
		}
	}
	if len(cluster) == 1 {
		result = append(result, point)
	} else {
		sort.Slice(cluster, func(a, b int) bool {
			return cluster[a].IntensityDifference > cluster[b].IntensityDifference
		})
		result = append(result, cluster[0])
	}
	return nmsr(radius, leftovers, result)
}
