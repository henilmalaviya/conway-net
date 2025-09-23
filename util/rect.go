package util

import (
	"math"

	"github.com/henilmalaviya/gol/grid"
)

func DiagonalLength(rect grid.Rectangle) float64 {
	width := float64(rect.Width())
	height := float64(rect.Height())

	w2 := width * width
	h2 := height * height

	return math.Sqrt(w2 + h2)
}
