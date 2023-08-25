package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 600
	screenHeight = 600
	boardSize    = 120
)

var (
	black = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	white = color.RGBA{R: 255, G: 255, B: 255, A: 255}
)

type Tile struct {
	on    bool
	x     int
	y     int
	image *ebiten.Image
}

type Board struct {
	size  int
	tiles map[*Tile]struct{}
	image *ebiten.Image
}

type Game struct {
	board *Board
}

func createTile(x, y int, tileSize int) *Tile {
	image := ebiten.NewImage(tileSize, tileSize)
	image.Fill(black)

	return &Tile{
		on:    false,
		image: image,
		x:     x,
		y:     y,
	}
}

func createBoard(size int, screenWidth, screenHeight int) *Board {
	tileSize := screenWidth / size
	image := ebiten.NewImage(screenWidth, screenHeight)
	tiles := make(map[*Tile]struct{})

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			tile := createTile(x, y, tileSize)
			tiles[tile] = struct{}{}
		}
	}

	return &Board{
		size:  size,
		tiles: tiles,
		image: image,
	}
}

func (t *Tile) click(x, y int) {
	tileX := x / t.image.Bounds().Dx()
	tileY := y / t.image.Bounds().Dy()

	if t.x == tileX && t.y == tileY {
		t.on = true
		t.image.Fill(white)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 300, 300
}

func (g *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		for tile := range g.board.tiles {
			tile.click(x, y)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.board.image, nil)

	tileSize := g.board.image.Bounds().Dx() / g.board.size

	for tile := range g.board.tiles {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(tileSize*tile.x), float64(tileSize*tile.y))
		screen.DrawImage(tile.image, op)
	}
}

func main() {
	board := createBoard(boardSize, screenWidth, screenHeight)

	g := &Game{
		board: board,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("John Conway's Game of Life")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
