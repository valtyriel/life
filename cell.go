package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Cell struct {
	drawable uint32

	IsAlive   bool
	AliveNext bool

	x int
	y int
}

func (c *Cell) Draw() {

	if !c.IsAlive {
		return
	}
	gl.BindVertexArray(c.drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
}

func (c *Cell) CheckState(cells [][]*Cell) {
	c.IsAlive = c.AliveNext
	c.AliveNext = c.IsAlive

	liveCount := c.LiveNeighbors(cells)
	if c.IsAlive {
		switch {
		// 1. Any live cell with fewer than two neighbors dies,
		// as if caused by underpopulation.
		case liveCount < 2:
			c.AliveNext = false
		// 2. Any live cell with two or three neightbors
		// lives on to the next generation.
		case liveCount == 2 || liveCount == 3:
			c.AliveNext = true
		// 3. Any live cell with more than three neighbors dies,
		// as if by overpopulation.
		case liveCount > 3:
			c.AliveNext = false
		}
		// 4. Any dead cell with exactly three live neighbors becomes
		// a live cell, as if by reproduction.
	} else if liveCount == 3 {
		c.AliveNext = true
	}
}

func (c *Cell) LiveNeighbors(cells [][]*Cell) int {
	var liveCount int
	add := func(x, y int) {
		// If we're at an edge, check the other side of the board.
		if x == len(cells) {
			x = 0
		} else if x == -1 {
			x = len(cells) - 1
		}
		if y == len(cells[x]) {
			y = 0
		} else if y == -1 {
			y = len(cells[x]) - 1
		}

		if cells[x][y].IsAlive {
			liveCount++
		}
	}

	add(c.x-1, c.y)   // left
	add(c.x+1, c.y)   // right
	add(c.x, c.y+1)   // up
	add(c.x, c.y-1)   // down
	add(c.x-1, c.y+1) // top-left
	add(c.x+1, c.y+1) // top-right
	add(c.x-1, c.y-1) // bottom-left
	add(c.x+1, c.y-1) // bottom-right

	return liveCount
}

func newCell(x, y int) *Cell {
	points := make([]float32, len(square), len(square))
	copy(points, square)

	for i, v := range points {
		var position, size float32
		switch i % 3 {
		case 0:
			size = 1.0 / float32(columns)
			position = float32(x) * size
		case 1:
			size = 1.0 / float32(columns)
			position = float32(y) * size
		default:
			continue
		}

		if v < 0 {
			points[i] = 2*position - 1
		} else {
			points[i] = 2*(position+size) - 1
		}
	}

	return &Cell{
		drawable: makeVao(points),
		x:        x,
		y:        y,
	}
}
