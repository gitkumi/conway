package main

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 600
	screenHeight = 600
	gridSize     = 120
)

var (
	black = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	white = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	red   = color.RGBA{R: 255, G: 0, B: 0, A: 255}

	gray50  = color.RGBA{0xf9, 0xfa, 0xfb, 0xff}
	gray100 = color.RGBA{0xf3, 0xf4, 0xf6, 0xff}
	gray200 = color.RGBA{0xe5, 0xe7, 0xeb, 0xff}
	gray300 = color.RGBA{0xd1, 0xd5, 0xdb, 0xff}
	gray400 = color.RGBA{0x9c, 0xa3, 0xaf, 0xff}
	gray500 = color.RGBA{0x6b, 0x72, 0x80, 0xff}
	gray600 = color.RGBA{0x4b, 0x55, 0x63, 0xff}
	gray700 = color.RGBA{0x37, 0x41, 0x51, 0xff}
	gray800 = color.RGBA{0x1f, 0x29, 0x37, 0xff}
	gray900 = color.RGBA{0x11, 0x18, 0x27, 0xff}
	gray950 = color.RGBA{0x03, 0x07, 0x12, 0xff}
)

type Cell struct {
	live       bool
	immunity   int
	currentAge int
	x          int
	y          int
	image      *ebiten.Image
	// Index of neighbors
	neighbors []int
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

func (c *Cell) kill() {
	c.live = false
	c.immunity = 0
	c.currentAge = 0
	c.image.Fill(white)
}

func (c *Cell) spawn(liveNeighbors, immunity int) {
	c.live = true
	c.immunity = immunity

	if c.immunity > 0 {
		color := ternary(liveNeighbors > 2, gray200, white)
		c.image.Fill(color)
	} else {
		c.image.Fill(gray200)
	}
}

func (c *Cell) age() {
	c.currentAge++

	if c.immunity > 0 {
		c.immunity--
		return
	}

	switch c.currentAge {
	case 1:
		c.image.Fill(gray400)
	case 2:
		c.image.Fill(gray500)
	case 3:
		c.image.Fill(gray600)
	case 4:
		c.image.Fill(gray700)
	case 5:
		c.image.Fill(gray800)
	case 6:
		c.image.Fill(gray300)
	case 7:
		c.image.Fill(gray200)
	case 8:
		c.image.Fill(gray100)
	case 9:
		c.image.Fill(gray50)
	default:
		c.image.Fill(gray100)
	}
}

func (c *Cell) initNeighbors(cells []*Cell) {
	c.neighbors = make([]int, 0, 8)
	row, col := c.y, c.x

	for rowIndex := row - 1; rowIndex <= row+1; rowIndex++ {
		for colIndex := col - 1; colIndex <= col+1; colIndex++ {
			// Self
			if rowIndex == row && colIndex == col {
				continue
			}

			wrappedRow := (rowIndex + gridSize) % gridSize
			wrappedCol := (colIndex + gridSize) % gridSize

			c.neighbors = append(c.neighbors, wrappedRow*gridSize+wrappedCol)
		}
	}
}

func (c *Cell) countLiveNeighbors(cells []*Cell) int {
	alive := 0
	for _, index := range c.neighbors {
		neighbor := cells[index]
		if neighbor.live {
			alive++
		}
	}

	return alive
}

func (c *Cell) update(liveNeighbors int) {
	if c.live {
		c.age()

		if c.immunity <= 0 && (liveNeighbors < 2 || liveNeighbors > 3) {
			c.kill()
		}
	} else {
		if liveNeighbors == 3 {
			c.spawn(liveNeighbors, 0)
		}
	}
}

func (g *Game) tick() {
	nextCells := make([]*Cell, len(g.grid.cells))

	for i, cell := range g.grid.cells {
		nextCell := *cell
		liveNeighbors := nextCell.countLiveNeighbors(g.grid.cells)
		nextCell.update(liveNeighbors)
		nextCells[i] = &nextCell
	}

	g.grid.cells = nextCells
}

func (g *Game) spawnCells() {
	var deadCells []*Cell

	for _, c := range g.grid.cells {
		if !c.live {
			deadCells = append(deadCells, c)
		}
	}

	for _, dc := range deadCells {
		r := rand.Intn(100) + 1
		immunity := rand.Intn(8)

		if r <= 1 {
			liveNeighbors := dc.countLiveNeighbors(g.grid.cells)
			dc.spawn(liveNeighbors, immunity)
		}
	}
}

func (g *Game) Update() error {
	chance := rand.Intn(5) + 3

	if g.grid.generation%chance == 0 {
		g.spawnCells()
	}

	g.tick()
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

func createGrid(size int, screenWidth, screenHeight int) *Grid {
	return &Grid{
		generation: 0,
		size:       size,
		image:      ebiten.NewImage(screenWidth, screenHeight),
		cells:      createCells(gridSize),
	}
}

func createCells(size int) []*Cell {
	cellSize := screenWidth / size
	var cells []*Cell

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			cell := createCell(x, y, cellSize)
			cells = append(cells, cell)
		}
	}

	for _, cell := range cells {
		cell.initNeighbors(cells)
	}

	return cells
}

func createCell(x, y, cellSize int) *Cell {
	image := ebiten.NewImage(cellSize, cellSize)
	image.Fill(white)

	return &Cell{
		live:       false,
		currentAge: 0,
		image:      image,
		x:          x,
		y:          y,
	}
}

func ternary[T any](cond bool, left, right T) T {
	if cond {
		return left
	}

	return right
}
