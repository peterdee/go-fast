package main

import (
	"math"
	"sort"
)

// combine clusters if possible
func combineClusters(
	cluster []Point,
	clusters [][]Point,
	radius int,
) [][]Point {
	combine := false
	lastCluster := clusters[len(clusters)-1]
	for _, currentClusterPoint := range cluster {
		for _, lastClusterPoint := range lastCluster {
			if (currentClusterPoint.X-lastClusterPoint.X) < radius &&
				int(math.Abs(float64(currentClusterPoint.Y)-float64(lastClusterPoint.Y))) < radius {
				combine = true
				break
			}
		}
		if combine {
			break
		}
	}
	if !combine {
		return append(clusters, cluster)
	}
	clusterEnd := clamp(len(clusters)-2, 0, len(clusters))
	newClusters := clusters[0:clusterEnd]
	lastCluster = append(lastCluster, cluster...)
	sort.SliceStable(
		lastCluster,
		func(i, j int) bool {
			if lastCluster[i].X != lastCluster[j].X {
				return lastCluster[i].X < lastCluster[j].X
			}
			return lastCluster[i].Y < lastCluster[j].Y
		},
	)
	return append(newClusters, lastCluster)
}

// Non-max suppression (recursive)
func nms(
	array []Point,
	radius int,
	previous Point,
	cluster []Point,
	clusters [][]Point,
	isSorted bool,
) []Point {
	if !isSorted && len(array) > 0 {
		sort.SliceStable(
			array,
			func(i, j int) bool {
				if array[i].X != array[j].X {
					return array[i].X < array[j].X
				}
				return array[i].Y < array[j].Y
			},
		)
		return nms(
			array,
			radius,
			previous,
			cluster,
			clusters,
			true,
		)
	}
	if len(array) == 0 {
		var updatedClusters [][]Point
		if len(clusters) > 0 {
			updatedClusters = combineClusters(cluster, clusters, radius)
		} else {
			updatedClusters = append(updatedClusters, cluster)
		}
		result := []Point{}
		for _, points := range updatedClusters {
			sort.SliceStable(
				points,
				func(i, j int) bool {
					return points[i].IntensityDifference > points[j].IntensityDifference
				},
			)
			result = append(result, points[0])
		}
		return result
	}
	current, rest := array[0], array[1:]
	if previous.IsEmpty {
		return nms(
			rest,
			radius,
			current,
			[]Point{current},
			clusters,
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
			clusters,
			true,
		)
	}
	if len(clusters) == 0 {
		return nms(
			rest,
			radius,
			current,
			[]Point{current},
			append(clusters, cluster),
			true,
		)
	}
	return nms(
		rest,
		radius,
		current,
		[]Point{current},
		combineClusters(cluster, clusters, radius),
		true,
	)
}
