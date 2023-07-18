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
	primarySortField byte,
) [][]Point {
	combine := false
	lastCluster := clusters[len(clusters)-1]
	for _, currentClusterPoint := range cluster {
		for _, lastClusterPoint := range lastCluster {
			primaryDifference := currentClusterPoint.X - lastClusterPoint.X
			secondaryDifference := int(
				math.Abs(float64(currentClusterPoint.Y) - float64(lastClusterPoint.Y)),
			)
			if primarySortField != 'x' {
				primaryDifference = currentClusterPoint.Y - lastClusterPoint.Y
				secondaryDifference = int(
					math.Abs(float64(currentClusterPoint.X) - float64(lastClusterPoint.X)),
				)
			}
			if primaryDifference < radius && secondaryDifference < radius {
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
			if primarySortField == 'x' {
				if lastCluster[i].X != lastCluster[j].X {
					return lastCluster[i].X < lastCluster[j].X
				}
				return lastCluster[i].Y < lastCluster[j].Y
			} else {
				if lastCluster[i].Y != lastCluster[j].Y {
					return lastCluster[i].Y < lastCluster[j].Y
				}
				return lastCluster[i].X < lastCluster[j].X
			}
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
	primarySortField byte,
) []Point {
	if !isSorted && len(array) > 0 {
		sort.SliceStable(
			array,
			func(i, j int) bool {
				if primarySortField == 'x' {
					if array[i].X != array[j].X {
						return array[i].X < array[j].X
					}
					return array[i].Y < array[j].Y
				} else {
					if array[i].Y != array[j].Y {
						return array[i].Y < array[j].Y
					}
					return array[i].X < array[j].X
				}
			},
		)
		return nms(
			array,
			radius,
			previous,
			cluster,
			clusters,
			true,
			primarySortField,
		)
	}
	if len(array) == 0 {
		var updatedClusters [][]Point
		if len(clusters) > 0 {
			updatedClusters = combineClusters(cluster, clusters, radius, primarySortField)
		} else {
			updatedClusters = append(updatedClusters, cluster)
		}
		result := []Point{}
		for _, points := range updatedClusters {
			sort.SliceStable(
				points,
				func(i, j int) bool {
					return points[i].IntensitySum > points[j].IntensitySum
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
			primarySortField,
		)
	}
	primaryDifference := current.X - previous.X
	secondaryDifference := int(math.Abs(float64(current.Y) - float64(previous.Y)))
	if primarySortField != 'x' {
		primaryDifference = current.Y - previous.Y
		secondaryDifference = int(math.Abs(float64(current.X) - float64(previous.X)))
	}
	if primaryDifference < radius && secondaryDifference < radius {
		return nms(
			rest,
			radius,
			current,
			append(cluster, current),
			clusters,
			true,
			primarySortField,
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
			primarySortField,
		)
	}
	return nms(
		rest,
		radius,
		current,
		[]Point{current},
		combineClusters(cluster, clusters, radius, primarySortField),
		true,
		primarySortField,
	)
}

// apply NMS recursively until the length of point array stops changing
func nmsRecursion(
	array []Point,
	radius int,
	prevLength int,
	isFirst bool,
) []Point {
	clusteredX := nms(
		array,
		radius,
		Point{
			IsEmpty: true,
		},
		[]Point{},
		[][]Point{},
		false,
		'x',
	)
	lenX := len(clusteredX)
	if !isFirst && lenX == prevLength {
		return clusteredX
	}
	clusteredY := nms(
		clusteredX,
		radius,
		Point{
			IsEmpty: true,
		},
		[]Point{},
		[][]Point{},
		false,
		'y',
	)
	lenY := len(clusteredY)
	if lenX == lenY {
		return clusteredY
	}
	return nmsRecursion(clusteredY, radius, lenY, false)
}
