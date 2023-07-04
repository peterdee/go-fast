package main

import (
	"math"
	"sort"
)

// Non-max suppression (recursive)
// TODO: improve NMS to handle edge cases
func nms(
	points []Point,
	radius int,
	previous Point,
	cluster []Point,
	selected []Point,
	isSorted bool,
) []Point {
	if !isSorted {
		sort.SliceStable(
			points,
			func(i, j int) bool {
				if points[i].X != points[j].X {
					return points[i].X < points[j].X
				}
				return points[i].Y < points[j].Y
			},
		)
		return nms(
			points,
			radius,
			previous,
			cluster,
			selected,
			true,
		)
	}
	if len(points) == 0 {
		if len(cluster) > 0 {
			if len(cluster) > 1 {
				sort.Slice(
					cluster,
					func(i, j int) bool {
						return cluster[i].IntensityDifference > cluster[j].IntensityDifference
					},
				)
			}
			selected = append(selected, cluster[0])
		}
		return selected
	}
	current, rest := points[0], points[1:]
	if previous.IsEmpty {
		return nms(
			rest,
			radius,
			current,
			[]Point{current},
			selected,
			true,
		)
	}
	if (current.X-previous.X) < radius &&
		int(math.Abs(float64(current.Y)-float64(previous.Y))) < radius {
		return nms(
			rest,
			radius,
			current,
			append(cluster, current),
			selected,
			true,
		)
	}
	sort.Slice(
		cluster,
		func(i, j int) bool {
			return cluster[i].IntensityDifference > cluster[j].IntensityDifference
		},
	)
	return nms(
		rest,
		radius,
		current,
		[]Point{current},
		append(selected, cluster[0]),
		true,
	)
}
