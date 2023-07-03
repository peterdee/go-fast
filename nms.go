package main

import (
	"fmt"
	"math"
	"sort"
)

// Non-max suppression (recursive)
func nms(
	points []Point,
	radius int,
	previous Point,
	cluster []Point,
	selected []Point,
	isSorted bool,
	iteration int,
	clusters [][]Point,
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
		fmt.Println("sorted", points)
		return nms(
			points,
			radius,
			previous,
			cluster,
			selected,
			true,
			iteration+1,
			clusters,
		)
	}
	if len(points) == 0 {
		// add point to selected only if cluster is not empty
		if len(cluster) > 0 {
			// sort only if there are several elements in the cluster
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
		fmt.Println("iterations", iteration)
		fmt.Println("clusters", append(clusters, cluster))
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
			iteration+1,
			clusters,
		)
	}
	if (current.X-previous.X) < radius &&
		int(math.Abs(float64(current.Y)-float64(previous.Y))) < radius {
		fmt.Println(math.Abs(float64(current.Y)-float64(previous.Y)), current.Y, previous.Y, cluster)
		x := append(cluster, current)
		return nms(
			rest,
			radius,
			current,
			x,
			selected,
			true,
			iteration+1,
			clusters,
		)
	}
	sort.Slice(
		cluster,
		func(i, j int) bool {
			return cluster[i].IntensityDifference > cluster[j].IntensityDifference
		},
	)
	x := append(selected, cluster[0])
	return nms(
		rest,
		radius,
		current,
		[]Point{current},
		x,
		true,
		iteration+1,
		append(clusters, cluster),
	)
}
