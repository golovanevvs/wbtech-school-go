package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x, y float64
}

func NewPoint(x, y float64) *Point {
	return &Point{x, y}
}

func (p Point) Distance(other Point) float64 {
	dx := p.x - other.x
	dy := p.y - other.y

	distance := math.Sqrt(dx*dx + dy*dy)

	return distance
}

func main() {
	in := bufio.NewReader(os.Stdin)

	points := make([]*Point, 2)

	for numPoint := range 2 {
	loop:
		for {
			fmt.Printf("Enter the coordinates of the point %d (x,y): ", numPoint+1)

			str, err := in.ReadString('\n')
			if err != nil {
				fmt.Printf("Read error: %v\n", err)
				continue
			}

			strs := strings.Split(str, ",")

			if len(strs) != 2 {
				fmt.Println("Enter a string in the format x1,y1")
				continue
			}

			values := make([]float64, 2)

			for i, v := range strs {
				values[i], err = strconv.ParseFloat(strings.TrimSpace(v), 64)
				if err != nil {
					fmt.Printf("Value error: %v\n", err)
					continue loop
				}
			}

			points[numPoint] = NewPoint(values[0], values[1])

			break
		}
	}

	distance := points[0].Distance(*points[1])

	fmt.Printf("Distance: %.2f\n", distance)
}
