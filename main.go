package main

import (
	"fmt"
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
	neighbors  []int
}

type Grid struct {
	size       int
	cells      []*Cell
	image      *ebiten.Image
	generation int
}

type Game struct {
	grid       *Grid
	oldestCell int
}

func createCell(generation, x, y int, image *ebiten.Image) *Cell {
	newImage := ebiten.NewImageFromImage(image)

	return &Cell{
		live:       false,
		currentAge: 0,
		image:      newImage,
		x:          x,
		y:          y,
	}
}

func createGrid(size int, screenWidth, screenHeight int) *Grid {
	cellSize := screenWidth / size
	var cells []*Cell

	cellImage := ebiten.NewImage(cellSize, cellSize)
	cellImage.Fill(white)

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			cell := createCell(0, x, y, cellImage)
			cells = append(cells, cell)
		}
	}

	for _, cell := range cells {
		cell.initNeighbors(cells)
	}

	return &Grid{
		size:       size,
		cells:      cells,
		image:      ebiten.NewImage(screenWidth, screenHeight),
		generation: 0,
	}
}

func (c *Cell) kill() {
	c.live = false
	c.immunity = 0
	c.currentAge = 0
	c.image.Fill(white)
}

func (c *Cell) spawn(cells []*Cell, immunity int) {
	c.live = true
	c.immunity = immunity

	if c.immunity > 0 {
		liveNeighborsCount := c.countLiveNeighbors(cells)
		color := ternary(liveNeighborsCount > 2, gray200, white)
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

	if c.currentAge == 1 {
		c.image.Fill(gray400)
	} else if c.currentAge == 2 {
		c.image.Fill(gray500)
	} else if c.currentAge == 3 {
		c.image.Fill(gray600)
	} else if c.currentAge == 4 {
		c.image.Fill(gray700)
	} else if c.currentAge == 5 {
		c.image.Fill(gray800)
	} else if c.currentAge == 6 {
		c.image.Fill(gray300)
	} else if c.currentAge == 7 {
		c.image.Fill(gray200)
	} else if c.currentAge == 8 {
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

func (g *Game) tick() {
	nextCells := make([]*Cell, len(g.grid.cells))

	for i, cell := range g.grid.cells {
		liveNeighborsCount := cell.countLiveNeighbors(g.grid.cells)
		nextCell := *cell

		if nextCell.live {
			nextCell.age()

			g.oldestCell = ternary(
				nextCell.currentAge > g.oldestCell,
				nextCell.currentAge,
				g.oldestCell,
			)

			if nextCell.immunity <= 0 && (liveNeighborsCount < 2 || liveNeighborsCount > 3) {
				nextCell.kill()
			}
		} else {
			if liveNeighborsCount == 3 {
				nextCell.spawn(g.grid.cells, 0)
			}
		}

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
			dc.spawn(g.grid.cells, immunity)
		}
	}
}

func (g *Game) Update() error {
	chance := rand.Intn(5) + 3

	if g.grid.generation%chance == 0 {
		g.spawnCells()
	}

	g.tick()

	fmt.Println("Generation: ", g.grid.generation)
	fmt.Println("Oldest: ", g.oldestCell)
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

func ternary[T any](cond bool, left, right T) T {
	if cond {
		return left
	}

	return right
}
