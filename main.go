package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 500
	screenHeight = 500
	gridSize     = 100
)

var (
	black = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	white = color.RGBA{R: 255, G: 255, B: 255, A: 255}
)

type Cell struct {
	live  bool
	x     int
	y     int
	image *ebiten.Image
}

type Grid struct {
	size       int
	cells      []*Cell
	image      *ebiten.Image
	generation int
}

type Game struct {
	grid *Grid
}

func createCell(x, y int, cellSize int) *Cell {
	image := ebiten.NewImage(cellSize, cellSize)
	image.Fill(white)

	return &Cell{
		live:  false,
		image: image,
		x:     x,
		y:     y,
	}
}

func createGrid(size int, screenWidth, screenHeight int) *Grid {
	cellSize := screenWidth / size
	image := ebiten.NewImage(screenWidth, screenHeight)
	var cells []*Cell

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			cell := createCell(x, y, cellSize)
			cells = append(cells, cell)
		}
	}

	return &Grid{
		size:       size,
		cells:      cells,
		image:      image,
		generation: 0,
	}
}

func (t *Cell) toggle() {
	t.live = !t.live

	if t.live {
		t.image.Fill(black)
	} else {
		t.image.Fill(white)
	}
}

func (g *Game) handleClick() {
	x, y := ebiten.CursorPosition()

	for _, cell := range g.grid.cells {
		cellY := y / cell.image.Bounds().Dy()
		cellX := x / cell.image.Bounds().Dx()

		if cell.x == cellX && cell.y == cellY {
			cell.toggle()
		}
	}
}

func (c *Cell) getLiveNeighborsCount(cells []*Cell) int {
	neighbors := make([]*Cell, 0, 8)
	row, col := c.y, c.x

	for r := row - 1; r <= row+1; r++ {
		for c := col - 1; c <= col+1; c++ {
			// Self
			if r == row && c == col {
				continue
			}
			if r >= 0 && r < gridSize && c >= 0 && c < gridSize {
				neighbors = append(neighbors, cells[r*gridSize+c])
			}
		}
	}

	alive := 0
	for _, neighbor := range neighbors {
		if neighbor.live {
			alive++
		}
	}

	return alive
}

func (g *Game) tick() {
	nextCells := make([]*Cell, len(g.grid.cells))

	for i, cell := range g.grid.cells {
		liveNeighborsCount := cell.getLiveNeighborsCount(g.grid.cells)
		nextCell := *cell

		if cell.live {
			if liveNeighborsCount < 2 || liveNeighborsCount > 3 {
				nextCell.toggle()
			}
		} else {
			if liveNeighborsCount == 3 {
				nextCell.toggle()
			}
		}

		nextCells[i] = &nextCell
	}

	g.grid.cells = nextCells
}

func (g *Game) birth() {
	var deadCells []*Cell

	for _, c := range g.grid.cells {
		if !c.live {
			deadCells = append(deadCells, c)
		}
	}

	deadRatio := len(deadCells) / len(g.grid.cells)

	for _, dc := range deadCells {
		rand := rand.Intn(100) + 1

		if rand <= 1+(5*deadRatio) {
			dc.toggle()
		}
	}
}

func (g *Game) death() {
	var aliveCells []*Cell

	for _, c := range g.grid.cells {
		if c.live {
			aliveCells = append(aliveCells, c)
		}
	}

	aliveRatio := len(aliveCells) / len(g.grid.cells)

	for _, ac := range aliveCells {
		rand := rand.Intn(100) + 1

		if rand <= 1+(5*aliveRatio) {
			ac.toggle()
		}
	}
}

func (g *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.handleClick()
	} else {
		if g.grid.generation%3 == 0 {
			g.birth()
			g.death()
		}
		g.tick()
	}

	fmt.Println("Generation: ", g.grid.generation)
	g.grid.generation++

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.grid.image, nil)
	cellSize := g.grid.image.Bounds().Dx() / g.grid.size

	for _, cell := range g.grid.cells {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(cellSize*cell.x), float64(cellSize*cell.y))
		screen.DrawImage(cell.image, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 300, 300
}

func main() {
	grid := createGrid(gridSize, screenWidth, screenHeight)

	g := &Game{
		grid: grid,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("John Conway's Game of Life")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
