package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 420
	screenHeight = 420
	boardSize    = 10
	tileSize     = 42
	tileGap      = 2
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

func createTile(x, y int) *Tile {
	image := ebiten.NewImage(tileSize, tileSize)
	image.Fill(black)

	return &Tile{
		on:    false,
		image: image,
		x:     x,
		y:     y,
	}
}

func createBoard(size int) *Board {
	image := ebiten.NewImage(screenWidth, screenHeight)
	tiles := make(map[*Tile]struct{})

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			tile := createTile(x, y)
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
	tileX := x / tileSize
	tileY := y / tileSize

	if t.x == tileX && t.y == tileY {
		t.on = true
		t.image.Fill(white)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 210, 210
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

	for tile := range g.board.tiles {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(tileSize*tile.x+tileGap*tile.x), float64(tileSize*tile.y+tileGap*tile.y))
		screen.DrawImage(tile.image, op)
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Tiles")
	board := createBoard(boardSize)

	g := &Game{
		board: board,
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
