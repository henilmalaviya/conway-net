package util

import (
	"github.com/henilmalaviya/gol/grid"
	"github.com/tidwall/gjson"
)

func GetBoundsFromData(data gjson.Result, key string) (grid.Rectangle, bool) {
	if bd := data.Get(key); bd.Exists() {

		if !bd.IsArray() {
			return grid.Rectangle{}, false
		}

		bdArray := bd.Array()
		if len(bdArray) != 2 {
			return grid.Rectangle{}, false
		}

		min := bdArray[0]
		max := bdArray[1]

		if !min.IsArray() || !max.IsArray() {
			return grid.Rectangle{}, false
		}

		minArray := min.Array()
		maxArray := max.Array()
		if len(minArray) != 2 || len(maxArray) != 2 {
			return grid.Rectangle{}, false
		}

		x1 := minArray[0].Int()
		y1 := minArray[1].Int()
		x2 := maxArray[0].Int()
		y2 := maxArray[1].Int()

		if x1 == 0 && y1 == 0 && x2 == 0 && y2 == 0 {
			return grid.Rectangle{}, false
		}
		return *grid.NewRectangle(int(x1), int(y1), int(x2), int(y2)), true
	}
	return grid.Rectangle{}, false
}

func GetCellsArrayFromData(data gjson.Result, key string) ([]grid.Cell, bool) {
	var cellsArray []grid.Cell
	data.Get(key).ForEach(func(_, value gjson.Result) bool {
		cell := value.Array()
		if len(cell) == 2 {
			X := int(cell[0].Int())
			Y := int(cell[1].Int())
			c := *grid.NewCellFromCords(X, Y)
			cellsArray = append(cellsArray, c)
		}
		return true
	})
	return cellsArray, len(cellsArray) > 0
}
